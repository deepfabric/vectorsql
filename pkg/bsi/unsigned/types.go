package unsigned

import "github.com/pilosa/pilosa/roaring"

type ubsi struct {
	bitSize int
	ms      []*roaring.Bitmap
}

const (
	bsiExistsBit = 0
	bsiOffsetBit = 1
)
