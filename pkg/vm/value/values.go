package value

import (
	"github.com/deepfabric/vectorsql/pkg/vm/value/dynamic"
	"github.com/deepfabric/vectorsql/pkg/vm/value/static"
)

func (a *Bool) ToValues() Values {
	return &static.Bools{Vs: []bool{bool(*a)}}
}

func (a *Int) ToValues() Values {
	return &static.Ints{Vs: []int64{int64(*a)}}
}

func (a *Float) ToValues() Values {
	return &static.Floats{Vs: []float64{float64(*a)}}
}

func (a *String) ToValues() Values {
	return &dynamic.Strings{Vs: []string{string(*a)}}
}

func (a *Int8) ToValues() Values {
	return &static.Int8s{Vs: []int8{int8(*a)}}
}

func (a *Int16) ToValues() Values {
	return &static.Int16s{Vs: []int16{int16(*a)}}
}

func (a *Int32) ToValues() Values {
	return &static.Int32s{Vs: []int32{int32(*a)}}
}

func (a *Int64) ToValues() Values {
	return &static.Int64s{Vs: []int64{int64(*a)}}
}

func (a *Uint8) ToValues() Values {
	return &static.Uint8s{Vs: []uint8{uint8(*a)}}
}

func (a *Uint16) ToValues() Values {
	return &static.Uint16s{Vs: []uint16{uint16(*a)}}
}

func (a *Uint32) ToValues() Values {
	return &static.Uint32s{Vs: []uint32{uint32(*a)}}
}

func (a *Uint64) ToValues() Values {
	return &static.Uint64s{Vs: []uint64{uint64(*a)}}
}

func (a *Float32) ToValues() Values {
	return &static.Float32s{Vs: []float32{float32(*a)}}
}

func (a *Float64) ToValues() Values {
	return &static.Float64s{Vs: []float64{float64(*a)}}
}

func (a *Timestamp) ToValues() Values {
	return &static.Timestamps{Vs: []int64{int64(*a)}}
}
