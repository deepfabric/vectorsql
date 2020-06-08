package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/thinkkv/pkg/engine"
	"github.com/deepfabric/vectorsql/pkg/bsi"
	"github.com/deepfabric/vectorsql/pkg/bsi/signed"
	"github.com/deepfabric/vectorsql/pkg/bsi/unsigned"
	"github.com/deepfabric/vectorsql/pkg/vm/container/relation"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/util/encoding"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
	lru "github.com/hashicorp/golang-lru"
)

func init() {
	gob.Register(MetaData{})
}

func New(mcpu int, db engine.DB, bc, rc *lru.Cache) *storage {
	return &storage{
		db:   db,
		bc:   bc,
		rc:   rc,
		mcpu: mcpu,
	}
}

func (s *storage) Close() error {
	defer s.db.Close()
	return s.db.Sync()
}

func (s *storage) NewRelation(id string, md MetaData) error {
	data, err := encoding.Encode(md)
	if err != nil {
		return err
	}
	bat, err := s.db.NewBatch()
	if err != nil {
		return err
	}
	if err := bat.Set(metaKey(id), data); err != nil {
		bat.Cancel()
		return err
	}
	return bat.Commit()
}

func (s *storage) Relation(id string) (relation.Relation, error) {
	if v, ok := s.rc.Get(id); ok {
		return v.(relation.Relation), nil
	}
	r, err := s.getRelation(id)
	if err == nil {
		s.rc.Add(id, r)
	}
	return r, err

}

func (r *index) Destroy() error {
	return r.db.Sync()
}

func (r *index) String() string {
	return r.id
}

func (r *index) IsEvent() bool {
	r.RLock()
	defer r.RUnlock()
	return r.isE
}

func (r *index) IdBitmap() (bsi.Bsi, error) {
	r.RLock()
	defer r.RUnlock()
	return getUbsi(ubsiKey(r.id, "id"), r.db, r.lc)
}

func (r *index) Eq(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.mp[attr]; !ok {
		return nil, fmt.Errorf("attribute '%s' not exist", attr)
	}
	switch v.ResolvedType().Oid {
	case types.T_uint8:
		mp, err := getBitmap(buKey(r.id, attr, value.MustBeUint8(v)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp, nil
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
		iv := value.MustBeInt8(v)
		if iv >= 0 {
			return getBitmap(buKey(r.id, attr, uint8(iv)), r.db, r.lc)
		}
		return getBitmap(biKey(r.id, attr, uint8(-iv)), r.db, r.lc)
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
	case types.T_int:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeInt(v)))
	case types.T_float:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeFloat64(v) * 100000))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeFloat32(v) * 1000))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Eq(int64(value.MustBeFloat64(v) * 100000))
	case types.T_string:
		mp, err := getBitmap(bsKey(r.id, attr, value.MustBeString(v)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp, nil
	}
	return nil, fmt.Errorf("unsupport type '%s' for Eq", v.ResolvedType())
}

func (r *index) Ne(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.mp[attr]; !ok {
		return nil, fmt.Errorf("attribute '%s' not exist", attr)
	}
	switch v.ResolvedType().Oid {
	case types.T_uint8:
		return r.uint8Ne(r.db, r.id, attr, value.MustBeUint8(v))
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
		iv := value.MustBeInt8(v)
		if iv >= 0 {
			mp, err := getBitmap(biPkey(r.id, attr), r.db, r.lc)
			if err != nil {
				return nil, err
			}
			mq, err := r.uint8Ne(r.db, r.id, attr, uint8(iv))
			if err != nil {
				return nil, err
			}
			mp.Or(mq)
			return mp, nil
		}
		mp, err := getBitmap(buPkey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		mq, err := r.int8Ne(r.db, r.id, attr, uint8(-iv))
		if err != nil {
			return nil, err
		}
		mp.Or(mq)
		return mp, nil
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
	case types.T_int:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeInt(v)))
	case types.T_float:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeFloat64(v) * 100000))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeFloat32(v) * 1000))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ne(int64(value.MustBeFloat64(v) * 100000))
	}
	return nil, fmt.Errorf("unsupport type '%s' for Eq", v.ResolvedType())
}

func (r *index) Lt(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.mp[attr]; !ok {
		return nil, fmt.Errorf("attribute '%s' not exist", attr)
	}
	switch v.ResolvedType().Oid {
	case types.T_uint8:
		return r.uint8Lt(r.db, r.id, attr, value.MustBeUint8(v))
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
		iv := value.MustBeInt8(v)
		if iv >= 0 {
			mp, err := getBitmap(biPkey(r.id, attr), r.db, r.lc)
			if err != nil {
				return nil, err
			}
			mq, err := r.uint8Lt(r.db, r.id, attr, uint8(iv))
			if err != nil {
				return nil, err
			}
			mp.Or(mq)
			return mp, nil
		}
		return r.int8Gt(r.db, r.id, attr, uint8(-iv))
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
	case types.T_int:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeInt(v)))
	case types.T_float:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeFloat64(v) * 100000))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeFloat32(v) * 1000))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Lt(int64(value.MustBeFloat64(v) * 100000))
	}
	return nil, fmt.Errorf("unsupport type '%s' for Eq", v.ResolvedType())
}

func (r *index) Le(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.mp[attr]; !ok {
		return nil, fmt.Errorf("attribute '%s' not exist", attr)
	}
	switch v.ResolvedType().Oid {
	case types.T_uint8:
		return r.uint8Le(r.db, r.id, attr, value.MustBeUint8(v))
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
		iv := value.MustBeInt8(v)
		if iv >= 0 {
			mp, err := getBitmap(biPkey(r.id, attr), r.db, r.lc)
			if err != nil {
				return nil, err
			}
			mq, err := r.uint8Le(r.db, r.id, attr, uint8(iv))
			if err != nil {
				return nil, err
			}
			mp.Or(mq)
			return mp, nil
		}
		return r.int8Ge(r.db, r.id, attr, uint8(-iv))
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
	case types.T_int:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeInt(v)))
	case types.T_float:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeFloat64(v) * 100000))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeFloat32(v) * 1000))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Le(int64(value.MustBeFloat64(v) * 100000))
	}
	return nil, fmt.Errorf("unsupport type '%s' for Eq", v.ResolvedType())
}

func (r *index) Gt(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.mp[attr]; !ok {
		return nil, fmt.Errorf("attribute '%s' not exist", attr)
	}
	switch v.ResolvedType().Oid {
	case types.T_uint8:
		return r.uint8Gt(r.db, r.id, attr, value.MustBeUint8(v))
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
		iv := value.MustBeInt8(v)
		if iv < 0 {
			mp, err := getBitmap(buPkey(r.id, attr), r.db, r.lc)
			if err != nil {
				return nil, err
			}
			mq, err := r.uint8Lt(r.db, r.id, attr, uint8(-iv))
			if err != nil {
				return nil, err
			}
			mp.Or(mq)
			return mp, nil
		}
		return r.uint8Gt(r.db, r.id, attr, uint8(iv))
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
	case types.T_int:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeInt(v)))
	case types.T_float:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeFloat64(v) * 100000))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeFloat32(v) * 1000))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Gt(int64(value.MustBeFloat64(v) * 100000))
	}
	return nil, fmt.Errorf("unsupport type '%s' for Eq", v.ResolvedType())
}

func (r *index) Ge(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.mp[attr]; !ok {
		return nil, fmt.Errorf("attribute '%s' not exist", attr)
	}
	switch v.ResolvedType().Oid {
	case types.T_uint8:
		return r.uint8Ge(r.db, r.id, attr, value.MustBeUint8(v))
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
		iv := value.MustBeInt8(v)
		if iv < 0 {
			mp, err := getBitmap(buPkey(r.id, attr), r.db, r.lc)
			if err != nil {
				return nil, err
			}
			mq, err := r.uint8Le(r.db, r.id, attr, uint8(-iv))
			if err != nil {
				return nil, err
			}
			mp.Or(mq)
			return mp, nil
		}
		return r.uint8Ge(r.db, r.id, attr, uint8(iv))
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
	case types.T_int:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeInt(v)))
	case types.T_float:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeFloat64(v) * 100000))
	case types.T_float32:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeFloat32(v) * 1000))
	case types.T_float64:
		mp, err := getBsi(bsiKey(r.id, attr), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		return mp.Ge(int64(value.MustBeFloat64(v) * 100000))
	}
	return nil, fmt.Errorf("unsupport type '%s' for Eq", v.ResolvedType())
}

func (r *index) AddTuplesByJson(ts []map[string]interface{}) error {
	if len(ts) == 0 {
		return nil
	}
	r.Lock()
	defer r.Unlock()
	bat, err := r.db.NewBatch()
	if err != nil {
		return err
	}
	smp := make(map[string]bsi.Bsi)
	bmp := make(map[string]*roaring.Bitmap)
	for _, t := range ts {
		if err := r.addTuple(t, bat, bmp, smp); err != nil {
			bat.Cancel()
			return err
		}
	}
	for k, mp := range smp {
		v, err := mp.Show()
		if err != nil {
			bat.Cancel()
			return err
		}
		if err := bat.Set([]byte(k), v); err != nil {
			bat.Cancel()
			return err
		}
	}
	for k, mp := range bmp {
		v, err := mp.MarshalBinary()
		if err != nil {
			bat.Cancel()
			return err
		}
		if err := bat.Set([]byte(k), v); err != nil {
			bat.Cancel()
			return err
		}
	}
	if err := bat.Commit(); err != nil {
		bat.Cancel()
		return err
	}
	if err := r.db.Sync(); err != nil {
		return err
	}
	return nil
}

func (s *storage) getRelation(id string) (*index, error) {
	var r index

	r.id = id
	r.db = s.db
	r.lc = s.bc
	r.mcpu = s.mcpu
	r.mp = make(map[string]int32)
	{
		v, err := s.db.Get(metaKey(id))
		switch {
		case err == nil:
			var md MetaData

			r.isE = md.IsEvent
			if err := encoding.Decode(v, &md); err != nil {
				return nil, err
			}
			for i := range md.Types {
				r.mp[md.Attrs[i]] = md.Types[i]
			}
		case err != nil:
			return nil, err
		}

	}
	return &r, nil
}

func (r *index) addTuple(mp map[string]interface{}, bat engine.Batch, bmp map[string]*roaring.Bitmap, smp map[string]bsi.Bsi) error {
	id, err := getSeq(mp) // row number
	if err != nil {
		return err
	}
	if r.isE {
		if err := getId(mp); err != nil {
			return err
		}
	}
	for attr, e := range mp {
		switch t := e.(type) {
		case uint8:
			{
				k := buKey(r.id, attr, t)
				mp, ok := bmp[k]
				if !ok {
					var err error
					if mp, err = getBitmap(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = roaring.New()
					}
					bmp[k] = mp
				}
				mp.Add(id)
			}
		case uint16:
			{
				k := ubsiKey(r.id, attr)
				mp, ok := smp[k]
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
				mp.Set(id, uint64(t))
			}
		case uint32:
			{
				k := ubsiKey(r.id, attr)
				mp, ok := smp[k]
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
				mp.Set(id, uint64(t))
			}
		case uint64:
			{
				k := ubsiKey(r.id, attr)
				mp, ok := smp[k]
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
				mp.Set(id, t)
			}
		case int8:
			{
				{
					k := buKey(r.id, attr, uint8(t))
					if t < 0 {
						ut := uint8(-t)
						k = biKey(r.id, attr, ut)
					}
					mp, ok := bmp[k]
					if !ok {
						var err error
						if mp, err = getBitmap(k, r.db, r.lc); err != nil {
							return err
						}
						if mp == nil {
							mp = roaring.New()
						}
						bmp[k] = mp
					}
					mp.Add(id)
				}
				{
					k := buPkey(r.id, attr)
					if t < 0 {
						k = biPkey(r.id, attr)
					}
					mp, ok := bmp[k]
					if !ok {
						var err error
						if mp, err = getBitmap(k, r.db, r.lc); err != nil {
							return err
						}
						if mp == nil {
							mp = roaring.New()
						}
						bmp[k] = mp
					}
					mp.Add(id)
				}
			}
		case int16:
			{
				k := bsiKey(r.id, attr)
				mp, ok := smp[k]
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
				mp.Set(id, int64(t))
			}
		case int32:
			{
				k := bsiKey(r.id, attr)
				mp, ok := smp[k]
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
				mp.Set(id, int64(t))
			}
		case int64:
			{
				k := bsiKey(r.id, attr)
				mp, ok := smp[k]
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
				mp.Set(id, int64(t))
			}
		case int:
			{
				k := bsiKey(r.id, attr)
				mp, ok := smp[k]
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
				mp.Set(id, int64(t))
			}
		case string:
			{
				k := bsKey(r.id, attr, t)
				mp, ok := bmp[k]
				if !ok {
					var err error
					if mp, err = getBitmap(k, r.db, r.lc); err != nil {
						return err
					}
					if mp == nil {
						mp = roaring.New()
					}
					bmp[k] = mp
				}
				mp.Add(id)
			}
		case float32:
			{
				v := int64(t * 1000)
				{
					k := bsiKey(r.id, attr)
					mp, ok := smp[k]
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
					mp.Set(id, v)
				}
			}
		case float64:
			{
				v := int64(t * 100000)
				{
					k := bsiKey(r.id, attr)
					mp, ok := smp[k]
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
					mp.Set(id, v)
				}
			}
		}
	}
	return nil
}

func getBitmap(id string, db engine.DB, lc *lru.Cache) (*roaring.Bitmap, error) {
	if v, ok := lc.Get(id); ok {
		return v.(*roaring.Bitmap), nil
	}
	v, err := db.Get([]byte(id))
	switch {
	case err == nil:
		mp := roaring.New()
		if err := mp.UnmarshalBinary(v); err != nil {
			return nil, err
		}
		lc.Add(id, mp)
		return mp, nil
	case err == engine.NotExist:
		return nil, nil
	default:
		return nil, err
	}
}

func getBsi(id string, db engine.DB, lc *lru.Cache) (bsi.Bsi, error) {
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
		lc.Add(id, mp)
		return mp, nil
	case err == engine.NotExist:
		return nil, nil
	default:
		return nil, err
	}
}

func getUbsi(id string, db engine.DB, lc *lru.Cache) (bsi.Bsi, error) {
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
		lc.Add(id, mp)
		return mp, nil
	case err == engine.NotExist:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *index) int8Ne(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, 256)
	for i := 0; i < 256; i++ {
		if uint8(i) == v {
			continue
		}
		mp, err := getBitmap(biKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) int8Lt(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, 256)
	for i, j := 0, int(v); i < j; i++ {
		mp, err := getBitmap(biKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) int8Le(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 256)
	for i, j := 0, int(v)+1; i < j; i++ {
		mp, err := getBitmap(biKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) int8Gt(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, 256)
	for i := int(v) + 1; i < 256; i++ {
		mp, err := getBitmap(biKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) int8Ge(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, 256)
	for i := int(v); i < 256; i++ {
		mp, err := getBitmap(biKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) uint8Ne(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, 256)
	for i := 0; i < 256; i++ {
		if uint8(i) == v {
			continue
		}
		mp, err := getBitmap(buKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) uint8Lt(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, 256)
	for i, j := 0, int(v); i < j; i++ {
		mp, err := getBitmap(buKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) uint8Le(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 256)
	for i, j := 0, int(v)+1; i < j; i++ {
		mp, err := getBitmap(buKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) uint8Gt(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, 256)
	for i := int(v) + 1; i < 256; i++ {
		mp, err := getBitmap(buKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func (r *index) uint8Ge(db engine.DB, id, attr string, v uint8) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, 256)
	for i := int(v); i < 256; i++ {
		mp, err := getBitmap(buKey(r.id, attr, uint8(i)), r.db, r.lc)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			ms = append(ms, mp)
		}
	}
	return roaring.ParOr(r.mcpu, ms...), nil
}

func getId(mp map[string]interface{}) error {
	v, ok := mp[ID]
	if !ok {
		return errors.New("id not set")
	}
	if _, ok := v.(uint32); !ok {
		return errors.New("id is not uint32")
	}
	return nil
}

func getSeq(mp map[string]interface{}) (uint32, error) {
	v, ok := mp[SEQ]
	if !ok {
		return 0, errors.New("seq not set")
	}
	if _, ok := v.(uint32); !ok {
		return 0, errors.New("seq is not uint32")
	}
	delete(mp, SEQ)
	return v.(uint32), nil
}

func metaKey(id string) []byte {
	var buf bytes.Buffer

	buf.WriteString("M.")
	buf.WriteString(id)
	return buf.Bytes()
}

func bnKey(id, attr string) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	return buf.String()
}

func bbKey(id, attr string, v bool) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	buf.WriteByte('.')
	if v {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	return buf.String()
}

func biPkey(id, attr string) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	buf.WriteString(".I")
	return buf.String()
}

func buPkey(id, attr string) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	buf.WriteString(".U")
	return buf.String()
}

func biKey(id, attr string, v uint8) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	buf.WriteString(".I.")
	buf.WriteByte(byte(v))
	return buf.String()
}

func buKey(id, attr string, v uint8) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	buf.WriteString(".U.")
	buf.WriteByte(byte(v))
	return buf.String()
}

func bsKey(id, attr string, v string) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteByte('.')
	buf.WriteString(attr)
	buf.WriteByte('.')
	buf.WriteString(v)
	return buf.String()
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
