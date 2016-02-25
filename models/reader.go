package models

import (
     "io"
)

// Reader reads a packfile from a binary string splitting it on objects
type Reader struct {
	// MaxObjectsLimit is the limit of objects to be load in the packfile, if
	// a packfile excess this number an error is throw, the default value
	// is defined by DefaultMaxObjectsLimit, usually the default limit is more
	// than enough to work with any repository, working extremly big repositories
	// where the number of object is bigger the memory can be exhausted.
	MaxObjectsLimit uint32

	// Format specifies if we are using ref-delta's or ofs-delta's, choosing the
	// correct format the memory usage is optimized
	// https://github.com/git/git/blob/8d530c4d64ffcc853889f7b385f554d53db375ed/Documentation/technical/protocol-capabilities.txt#L154
//	Format Format

//	R       *trackingReader
	S       ObjectStorage
	offsets map[int64]Hash
}

// A PackReader implements access to Git pack files and indexes.
type PackReader struct {
	version   int
	pack, idx *io.SectionReader

	// idxFanout[i] is the number of objects whose first byte
	// is <= i.
	idxFanout [256]uint32
}

/*
const idxHeaderSize = 4 + 4 + 256*4

func (pk *PackReader) checkIdxMagic(idx *io.SectionReader) (err error) {
	var buf [idxHeaderSize]byte
	_, err = idx.ReadAt(buf[:], 0)
	if err != nil {
		return
	}
	magic := [4]byte{buf[0], buf[1], buf[2], buf[3]}
	if magic != ([4]byte{'\xff', 't', 'O', 'c'}) {
		return errBadIdxMagic
	}
	for i := range pk.idxFanout {
		pk.idxFanout[i] = binary.BigEndian.Uint32(buf[8+4*i:])
	}
	return nil
}

// Read reads the objects and stores it at the ObjectStorage
func (r *Reader) Read(s ObjectStorage) (int64, error) {
	r.S = s
	if err := r.validateHeader(); err != nil {
		if err == io.EOF {
			return -1, EmptyRepositoryErr
		}

		return -1, err
	}

	version, err := r.readInt32()
	if err != nil {
		return -1, err
	}

	if version > VersionSupported {
		return -1, UnsupportedVersionErr
	}

	count, err := r.readInt32()
	if err != nil {
		return -1, err
	}

	if count > r.MaxObjectsLimit {
		return -1, MaxObjectsLimitReachedErr
	}

	return r.R.Position, r.ReadObjects(count)
}

func (r *Reader) validateHeader() error {
	var header = make([]byte, 4)
	if _, err := io.ReadFull(r.R, header); err != nil {
		return err
	}
//    fmt.Printf("########  validateHeader() header: %v ########\r\n", string(header))
	if !bytes.Equal(header, []byte{'P', 'A', 'C', 'K'}) {
		return MalformedPackfileErr
	}

	return nil
}

func (r *Reader) readInt32() (uint32, error) {
	var value uint32
	if err := binary.Read(r.R, binary.BigEndian, &value); err != nil {
		return 0, err
	}
//    fmt.Printf("########  readInt32() value: %v ########\r\n", value)
	return value, nil
}

func (r *Reader) ReadObjects(count uint32) error {
	// This code has 50-80 µs of overhead per object not counting zlib inflation.
	// Together with zlib inflation, it's 400-410 µs for small objects.
	// That's 1 sec for ~2450 objects, ~4.20 MB, or ~250 ms per MB,
	// of which 12-20 % is _not_ zlib inflation (ie. is our code).
	for i := 0; i < int(count); i++ {
		start := r.R.Position
//        fmt.Printf("###### ReadObjects start :%v######\r\n", start)
		obj, err := r.NewRAWObject()
		if err != nil && err != io.EOF {
			return err
		}

		if r.Format == UnknownFormat || r.Format == OFSDeltaFormat {
			r.offsets[start] = obj.Hash()
		}

		r.S.Set(obj)
		if err == io.EOF {
			break
		}
	}

	return nil
}

func (r *Reader) NewRAWObject() (Object, error) {
	raw := r.S.New()
	var steps int64

	var buf [1]byte
	if _, err := r.R.Read(buf[:]); err != nil {
		return nil, err
	}
//    fmt.Printf("###### NewRAWObject() buf :%v######\r\n",buf)

	typ := ObjectType((buf[0] >> 4) & 7)
	size := int64(buf[0] & 15)
	steps++ // byte we just read to get `o.typ` and `o.size`

	var shift uint = 4
	for buf[0]&0x80 == 0x80 {
		if _, err := r.R.Read(buf[:]); err != nil {
			return nil, err
		}

		size += int64(buf[0]&0x7f) << shift
		steps++ // byte we just read to update `o.size`
		shift += 7
	}

	raw.SetType(typ)
	raw.SetSize(size)

//    fmt.Printf("###### NewRAWObject() typ :%v######\r\n",typ)
//    fmt.Printf("###### NewRAWObject() size :%v######\r\n",size)

	var err error
	switch raw.Type() {
	case REFDeltaObject:
		err = r.readREFDelta(raw)
	case OFSDeltaObject:
		err = r.readOFSDelta(raw, steps)
	case CommitObject, TreeObject, BlobObject, TagObject:
		err = r.readObject(raw)
	default:
		err = InvalidObjectErr.n("tag %q", raw.Type)
	}

	return raw, err
}

func (r *Reader) readREFDelta(raw Object) error {
	var ref Hash
	if _, err := io.ReadFull(r.R, ref[:]); err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err := r.inflate(buf); err != nil {
		return err
	}

	referenced, ok := r.S.Get(ref)
	if !ok {
		return ObjectNotFoundErr.n("%s", ref)
	}

	d, _ := ioutil.ReadAll(referenced.Reader())
	patched := patchDelta(d, buf.Bytes())
	if patched == nil {
		return PatchingErr.n("hash %q", ref)
	}

	raw.SetType(referenced.Type())
	raw.SetSize(int64(len(patched)))
	raw.Writer().Write(patched)

	return nil
}

func (r *Reader) readOFSDelta(raw Object, steps int64) error {
	start := r.R.Position
	offset, err := decodeOffset(r.R, steps)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err := r.inflate(buf); err != nil {
		return err
	}

	ref, ok := r.offsets[start+offset]
	if !ok {
		return PackEntryNotFoundErr.n("offset %d", start+offset)
	}

	referenced, _ := r.S.Get(ref)
	d, _ := ioutil.ReadAll(referenced.Reader())
	patched := patchDelta(d, buf.Bytes())
	if patched == nil {
		return PatchingErr.n("hash %q", ref)
	}

	raw.SetType(referenced.Type())
	raw.SetSize(int64(len(patched)))
	raw.Writer().Write(patched)

	return nil
}

func (r *Reader) readObject(raw Object) error {
	return r.inflate(raw.Writer())
}

func (r *Reader) inflate(w io.Writer) error {
	zr, err := zlib.NewReader(r.R)
	if err != nil {
		if err == zlib.ErrHeader {
			return zlib.ErrHeader
		}

		return ZLibErr.n("%s", err)
	}
//    fmt.Printf("###### inflate zr:%v######\r\n", zr)
	defer zr.Close()
     
	_, err = io.Copy(w, zr)
//    fmt.Printf("###### inflate w:%v######\r\n", w)
	return err
}

*/