package models

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	_"fmt"
)

// Hash SHA1 hased content
type Hash [20]byte

// ComputeHash compute the hash for a given ObjectType and content
func ComputeHash(t ObjType, content []byte) Hash {
	h := t.Bytes()
	h = append(h, ' ')
	h = strconv.AppendInt(h, int64(len(content)), 10)
	h = append(h, 0)
	h = append(h, content...)

	return Hash(sha1.Sum(h))
}

// NewHash return a new Hash from a hexadecimal hash representation
func NewHash(s string) Hash {
//	fmt.Printf("######   NewHash s: %v ########\r\n", s)
	b, _ := hex.DecodeString(s)
//	fmt.Printf("######   NewHash b: %v ########\r\n", b)
	var h Hash
	copy(h[:], b)

	return h
}

func (h Hash) IsZero() bool {
	var empty Hash
	return h == empty
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}
