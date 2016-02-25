package models

import (
	"bytes"
	"io"
)

// ObjectType internal object type's
type ObjType int8

// Object is a generic representation of any git object
type Object interface {
	Type() ObjType
	SetType(ObjType)
	Size() int64
	SetSize(int64)
	Hash() Hash
	Reader() io.Reader
	Writer() io.Writer
}

func (t ObjType) Bytes() []byte {
	return []byte(t.String())
}

const (
	CommitObject   ObjType = 1
	TreeObject     ObjType = 2
	BlobObject     ObjType = 3
	TagObject      ObjType = 4
	OFSDeltaObject ObjType = 6
	REFDeltaObject ObjType = 7
)

func (t ObjType) String() string {
	switch t {
	case CommitObject:
		return "commit"
	case TreeObject:
		return "tree"
	case BlobObject:
		return "blob"
	default:
		return "-"
	}
}

type RAWObject struct {
	b []byte
	t ObjType
	s int64
}
/*
// A loose object consists of
// <type> <size>\x00
// where type is "blob", "tree" or "commit"
func readLoose(r io.ReadCloser) (t ObjType, s []byte, err error) {
	// read compressed data.
	zr, err := zlib.NewReader(r)
	if err != nil {
		r.Close()
		return
	}
	s, err = ioutil.ReadAll(zr)
	r.Close()
    fmt.Printf("####### readLoose r:%v #######\r\n", r)
    fmt.Printf("####### readLoose zr:%v #######\r\n", zr)
    fmt.Printf("####### readLoose s:%v #######\r\n", s)
	hdr := s
	if len(hdr) > 32 {
		hdr = hdr[:32] // 32 bytes are enough for a 20-digit size.
	}
	if len(hdr) < 4 {
		err = errCorruptedObjectHeader
		return
	}
	sp := bytes.IndexByte(hdr, ' ')
	nul := bytes.IndexByte(hdr, 0)
	switch s[0] {
	case 'b':
		t = BLOB
	case 'c':
		t = COMMIT
	case 't':
		t = TREE
	}
	if sp < 0 || !matchType(t, hdr[:sp]) {
		err = errInvalidType(string(hdr[:sp]))
		return
	}
	if nul < 0 {
		err = errCorruptedObjectHeader
		return
	}
	sz, err := strconv.ParseUint(string(hdr[sp+1:nul]), 10, 64)
	if err != nil {
		err = errCorruptedObjectHeader
		return
	}
	s = s[nul+1:]
	if uint64(len(s)) != sz {
		err = errObjectSizeMismatch
	}
    fmt.Printf("####### readLoose t:%v #######\r\n", t)
	return t, s, err
}

func readObject(t ObjType, data []byte) (Object, error) {
	switch t {
	case BLOB:
		o := Blob{Data: data}
		o.Hash = rehash(o)
		return o, nil
	case TREE:
		o, err := parseTree(data)
		o.Hash = rehash(o)
		return o, err
	case COMMIT:
		o, err := parseCommit(data)
		o.Hash = rehash(o)
		return o, err
	}
	panic(errInvalidType(t.String()))
}

// ParseLoose reads a loose object as stored in the objects/
// subdirectory of a git repository.
func ParseLoose(r io.ReadCloser) (Object, error) {
	t, data, err := readLoose(r)
	if err != nil {
		return nil, err
	}
	return readObject(t, data)
}
*/
func (o *RAWObject) Type() ObjType     { return o.t }
func (o *RAWObject) SetType(t ObjType) { o.t = t }
func (o *RAWObject) Size() int64          { return o.s }
func (o *RAWObject) SetSize(s int64)      { o.s = s }
func (o *RAWObject) Reader() io.Reader    { return bytes.NewBuffer(o.b) }
func (o *RAWObject) Hash() Hash           { return ComputeHash(o.t, o.b) }
func (o *RAWObject) Writer() io.Writer    { return o }
func (o *RAWObject) Write(p []byte) (n int, err error) {
	o.b = append(o.b, p...)
	return len(p), nil
}

// ObjectStorage generic storage of objects
type ObjectStorage interface {
	New() Object
	Set(Object) Hash
	Get(Hash) (Object, bool)
}
