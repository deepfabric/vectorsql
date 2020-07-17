package static

import "github.com/pilosa/pilosa/roaring"

type Ints struct {
	Vs []int64
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Int8s struct {
	Vs []int8
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Int16s struct {
	Vs []int16
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Int32s struct {
	Vs []int32
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Int64s struct {
	Vs []int64
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Uint8s struct {
	Vs []uint8
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Uint16s struct {
	Vs []uint16
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Uint32s struct {
	Vs []uint32
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Uint64s struct {
	Vs []uint64
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Bools struct {
	Vs []bool
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Timestamps struct {
	Vs []int64
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Floats struct {
	Vs []float64
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Float32s struct {
	Vs []float32
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}

type Float64s struct {
	Vs []float64
	Is []uint64
	Np *roaring.Bitmap // null
	Dp *roaring.Bitmap // null
}
