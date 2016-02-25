package models

import (
	_"bytes"
	_"io"
)

type RAWObjectStorage struct {
	Objects map[Hash]Object
	Commits map[Hash]Object
	Trees   map[Hash]Object
	Blobs   map[Hash]Object
}

func NewRAWObjectStorage() *RAWObjectStorage {
	return &RAWObjectStorage{
		Objects: make(map[Hash]Object, 0),
		Commits: make(map[Hash]Object, 0),
		Trees:   make(map[Hash]Object, 0),
		Blobs:   make(map[Hash]Object, 0),
	}
}

func (o *RAWObjectStorage) New() Object {
	return &RAWObject{}
}

func (o *RAWObjectStorage) Set(obj Object) Hash {
	h := obj.Hash()
	o.Objects[h] = obj

	switch obj.Type() {
	case CommitObject:
		o.Commits[h] = o.Objects[h]
	case TreeObject:
		o.Trees[h] = o.Objects[h]
	case BlobObject:
		o.Blobs[h] = o.Objects[h]
	}

	return h
}

func (o *RAWObjectStorage) Get(h Hash) (Object, bool) {
	obj, ok := o.Objects[h]

	return obj, ok
}

