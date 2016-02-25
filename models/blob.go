package models

import (
	_"log"
	"io"
)

// Blob is used to store file data - it is generally a file.
type Blob struct {
	Hash Hash
	Size int64

	obj Object
}

// Decode transform an core.Object into a Blob struct
func (b *Blob) Decode(o Object) error {
	b.Hash = o.Hash()
	b.Size = o.Size()
	b.obj = o

	return nil
}

// Reader returns a reader allow the access to the content of the blob
func (b *Blob) Reader() io.Reader {
	return b.obj.Reader()
}