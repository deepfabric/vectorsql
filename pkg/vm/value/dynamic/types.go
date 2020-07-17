package dynamic

import "github.com/pilosa/pilosa/roaring"

type Strings struct {
	Vs []string
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}
