package index

import (
	"bytes"
	"fmt"

	"github.com/deepfabric/thinkkv/pkg/engine"
	"github.com/deepfabric/vectorsql/pkg/bsi"
	"github.com/deepfabric/vectorsql/pkg/bsi/signed"
	"github.com/deepfabric/vectorsql/pkg/bsi/unsigned"
	"github.com/deepfabric/vectorsql/pkg/storage/cache"
	"github.com/deepfabric/vectorsql/pkg/storage/metadata"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
	"github.com/pilosa/pilosa/roaring"
)

func New(isE bool, id string, db engine.DB, lc cache.Cache, attrs []metadata.Attribute) *index {
	return &index{
		id:    id,
		db:    db,
		lc:    lc,
		isE:   isE,
		attrs: attrs,
	}
}

func (r *index) IdMap() (bsi.Bsi, error) {
	return getUbsi(ubsiKey(r.id, ID), r.db, r.lc)
}

func (r *index) Eq(attr string, v value.Value) (*roaring.Bitmap, error) {
	switch v.ResolvedType() {
	case types.T_uint8:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(uint64(value.MustBeUint8(v)))
	case types.T_uint16:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(uint64(value.MustBeUint16(v)))
	case types.T_uint32:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(uint64(value.MustBeUint32(v)))
	case types.T_uint64:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(uint64(value.MustBeUint64(v)))
	case types.T_int8:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeInt8(v)))
	case types.T_int16:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeInt16(v)))
	case types.T_int32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeInt32(v)))
	case types.T_int64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeInt64(v)))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeFloat32(v) * Ffrac))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeFloat64(v) * Dfrac))
	case types.T_timestamp:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(value.MustBeTimestamp(v).Unix())
	}
	return nil, fmt.Errorf("unsupport type '%s' for Eq", v.ResolvedType())
}

func (r *index) Ne(attr string, v value.Value) (*roaring.Bitmap, error) {
	switch v.ResolvedType() {
	case types.T_uint8:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(uint64(value.MustBeUint8(v)))
	case types.T_uint16:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(uint64(value.MustBeUint16(v)))
	case types.T_uint32:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(uint64(value.MustBeUint32(v)))
	case types.T_uint64:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(uint64(value.MustBeUint64(v)))
	case types.T_int8:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeInt8(v)))
	case types.T_int16:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeInt16(v)))
	case types.T_int32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeInt32(v)))
	case types.T_int64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeInt64(v)))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeFloat32(v) * Ffrac))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeFloat64(v) * Dfrac))
	case types.T_timestamp:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(value.MustBeTimestamp(v).Unix())
	}
	return nil, fmt.Errorf("unsupport type '%s' for Ne", v.ResolvedType())
}

func (r *index) Lt(attr string, v value.Value) (*roaring.Bitmap, error) {
	switch v.ResolvedType() {
	case types.T_uint8:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(uint64(value.MustBeUint8(v)))
	case types.T_uint16:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(uint64(value.MustBeUint16(v)))
	case types.T_uint32:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(uint64(value.MustBeUint32(v)))
	case types.T_uint64:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(uint64(value.MustBeUint64(v)))
	case types.T_int8:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeInt8(v)))
	case types.T_int16:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeInt16(v)))
	case types.T_int32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeInt32(v)))
	case types.T_int64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeInt64(v)))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeFloat32(v) * Ffrac))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeFloat64(v) * Dfrac))
	case types.T_timestamp:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(value.MustBeTimestamp(v).Unix())
	}
	return nil, fmt.Errorf("unsupport type '%s' for Lt", v.ResolvedType())
}

func (r *index) Le(attr string, v value.Value) (*roaring.Bitmap, error) {
	switch v.ResolvedType() {
	case types.T_uint8:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(uint64(value.MustBeUint8(v)))
	case types.T_uint16:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(uint64(value.MustBeUint16(v)))
	case types.T_uint32:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(uint64(value.MustBeUint32(v)))
	case types.T_uint64:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(uint64(value.MustBeUint64(v)))
	case types.T_int8:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeInt8(v)))
	case types.T_int16:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeInt16(v)))
	case types.T_int32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeInt32(v)))
	case types.T_int64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeInt64(v)))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeFloat32(v) * Ffrac))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeFloat64(v) * Dfrac))
	case types.T_timestamp:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(value.MustBeTimestamp(v).Unix())
	}
	return nil, fmt.Errorf("unsupport type '%s' for Le", v.ResolvedType())
}

func (r *index) Gt(attr string, v value.Value) (*roaring.Bitmap, error) {
	switch v.ResolvedType() {
	case types.T_uint8:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(uint64(value.MustBeUint8(v)))
	case types.T_uint16:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(uint64(value.MustBeUint16(v)))
	case types.T_uint32:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(uint64(value.MustBeUint32(v)))
	case types.T_uint64:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(uint64(value.MustBeUint64(v)))
	case types.T_int8:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeInt8(v)))
	case types.T_int16:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeInt16(v)))
	case types.T_int32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeInt32(v)))
	case types.T_int64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeInt64(v)))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeFloat32(v) * Ffrac))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeFloat64(v) * Dfrac))
	case types.T_timestamp:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(value.MustBeTimestamp(v).Unix())
	}
	return nil, fmt.Errorf("unsupport type '%s' for Gt", v.ResolvedType())
}

func (r *index) Ge(attr string, v value.Value) (*roaring.Bitmap, error) {
	switch v.ResolvedType() {
	case types.T_uint8:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(uint64(value.MustBeUint8(v)))
	case types.T_uint16:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(uint64(value.MustBeUint16(v)))
	case types.T_uint32:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(uint64(value.MustBeUint32(v)))
	case types.T_uint64:
		mp, err := getUbsi(ubsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(uint64(value.MustBeUint64(v)))
	case types.T_int8:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeInt8(v)))
	case types.T_int16:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeInt16(v)))
	case types.T_int32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeInt32(v)))
	case types.T_int64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeInt64(v)))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeFloat32(v) * Ffrac))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeFloat64(v) * Dfrac))
	case types.T_timestamp:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(value.MustBeTimestamp(v).Unix())
	}
	return nil, fmt.Errorf("unsupport type '%s' for Ge", v.ResolvedType())
}

func (r *index) AddTuples(ts []interface{}) error {
	var seqs []uint64

	seqs = ts[0].([]uint64)
	smp := make(map[string]bsi.Bsi)
	for i, j := 1, len(ts); i < j; i++ {
		if err := r.addTuple(seqs, r.attrs[i], ts[i], smp); err != nil {
			return err
		}
	}
	for k, mp := range smp {
		v, err := mp.Show()
		if err != nil {
			return err
		}
		if err := r.db.Set([]byte(k), v); err != nil {
			return err
		}
	}
	return nil
}

func (r *index) addTuple(seqs []uint64, attr metadata.Attribute, t interface{}, smp map[string]bsi.Bsi) error {
	switch attr.Type {
	case types.T_timestamp:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := bsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getBsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = signed.New(64)
					}
					smp[k] = mp
				}
			}
			vs := t.([]int64)
			for i, v := range vs {
				mp.Set(seqs[i], v)
			}
		}
	case types.T_uint8:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := ubsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getUbsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = unsigned.New(8)
					}
					smp[k] = mp
				}
			}
			vs := t.([]uint8)
			for i, v := range vs {
				mp.Set(seqs[i], uint64(v))
			}
		}
	case types.T_uint16:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := ubsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getUbsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = unsigned.New(16)
					}
					smp[k] = mp
				}
			}
			vs := t.([]uint16)
			for i, v := range vs {
				mp.Set(seqs[i], uint64(v))
			}
		}
	case types.T_uint32:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := ubsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getUbsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = unsigned.New(32)
					}
					smp[k] = mp
				}
			}
			vs := t.([]uint32)
			for i, v := range vs {
				mp.Set(seqs[i], uint64(v))
			}
		}
	case types.T_uint64:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := ubsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getUbsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = unsigned.New(64)
					}
					smp[k] = mp
				}
			}
			vs := t.([]uint64)
			for i, v := range vs {
				mp.Set(seqs[i], uint64(v))
			}
		}
	case types.T_int8:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := bsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getBsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = signed.New(8)
					}
					smp[k] = mp
				}
			}
			vs := t.([]int8)
			for i, v := range vs {
				mp.Set(seqs[i], int64(v))
			}
		}
	case types.T_int16:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := bsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getBsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = signed.New(16)
					}
					smp[k] = mp
				}
			}
			vs := t.([]int16)
			for i, v := range vs {
				mp.Set(seqs[i], int64(v))
			}
		}
	case types.T_int32:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := bsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getBsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = signed.New(32)
					}
					smp[k] = mp
				}
			}
			vs := t.([]int32)
			for i, v := range vs {
				mp.Set(seqs[i], int64(v))
			}
		}
	case types.T_int64:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := bsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getBsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = signed.New(64)
					}
					smp[k] = mp
				}
			}
			vs := t.([]int64)
			for i, v := range vs {
				mp.Set(seqs[i], v)
			}
		}
	case types.T_float32:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := bsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getBsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = signed.New(32)
					}
					smp[k] = mp
				}
			}
			vs := t.([]float32)
			for i, v := range vs {
				mp.Set(seqs[i], int64(v*Ffrac))
			}
		}
	case types.T_float64:
		{
			var mp bsi.Bsi

			{
				var ok bool
				k := bsiKey(r.id, attr.Name)
				mp, ok = smp[k]
				if !ok {
					var err error
					if mp, err = getBsi(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = signed.New(64)
					}
					smp[k] = mp
				}
			}
			vs := t.([]float64)
			for i, v := range vs {
				mp.Set(seqs[i], int64(v*Dfrac))
			}
		}
	}
	return nil
}

func getBsi(id string, db engine.DB, lc cache.Cache) (bsi.Bsi, error) {
	if v, ok := lc.Get(id); ok {
		return v.(bsi.Bsi), nil
	}
	v, err := db.Get([]byte(id))
	switch {
	case err == nil:
		mp := signed.New(0)
		if err := mp.Read(v); err != nil {
			return nil, err
		}
		lc.Set(id, mp, v)
		return mp, nil
	case err == engine.NotExist:
		return nil, nil
	default:
		return nil, err
	}
}

func getUbsi(id string, db engine.DB, lc cache.Cache) (bsi.Bsi, error) {
	if v, ok := lc.Get(id); ok {
		return v.(bsi.Bsi), nil
	}
	v, err := db.Get([]byte(id))
	switch {
	case err == nil:
		mp := unsigned.New(0)
		if err := mp.Read(v); err != nil {
			return nil, err
		}
		lc.Set(id, mp, v)
		return mp, nil
	case err == engine.NotExist:
		return nil, nil
	default:
		return nil, err
	}
}

func getBitmap(id string, db engine.DB, lc cache.Cache) (*roaring.Bitmap, error) {
	if v, ok := lc.Get(id); ok {
		return v.(*roaring.Bitmap), nil
	}
	v, err := db.Get([]byte(id))
	switch {
	case err == nil:
		mp := roaring.NewBitmap()
		if err := mp.UnmarshalBinary(v); err != nil {
			return nil, err
		}
		lc.Set(id, mp, v)
		return mp, nil
	case err == engine.NotExist:
		return nil, nil
	default:
		return nil, err
	}
}

func bsiKey(id, attr string) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	buf.WriteString(".I")
	return buf.String()
}

func ubsiKey(id, attr string) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	buf.WriteString(".U")
	return buf.String()
}

func show(mp *roaring.Bitmap) ([]byte, error) {
	var buf bytes.Buffer

	if _, err := mp.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
