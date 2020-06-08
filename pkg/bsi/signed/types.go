package signed

import "github.com/pilosa/pilosa/roaring"

type bsi struct {
	bitSize int
	ms      []*roaring.Bitmap
}

const (
	bsiExistsBit = 0
	bsiSignBit   = 1
	bsiOffsetBit = 2
)
