package storage

import (
	"fmt"

	"github.com/deepfabric/thinkkv/pkg/engine"
	"github.com/deepfabric/vectorsql/pkg/bsi"
	"github.com/deepfabric/vectorsql/pkg/lru"
	"github.com/deepfabric/vectorsql/pkg/storage/cache"
	"github.com/deepfabric/vectorsql/pkg/storage/index"
	"github.com/deepfabric/vectorsql/pkg/storage/metadata"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/util/encoding"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
	"github.com/pilosa/pilosa/roaring"
)

func New(db engine.DB, rc lru.LRU, lc cache.Cache) *storage {
	return &storage{
		db: db,
		lc: lc,
		rc: rc,
	}
}

func (s *storage) Close() error {
	defer s.db.Close()
	return s.db.Sync()
}

func (s *storage) Relation(id string) (Relation, error) {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.rc.Get(id); ok {
		return v.(Relation), nil
	}
	r := new(relation)
	{
		v, err := s.db.Get(metadata.Mkey(id))
		if err != nil {
			return nil, err
		}
		if err := encoding.Decode(v, &r.md); err != nil {
			return nil, err
		}
	}
	r.idx = index.New(r.md.IsE, id, s.db, s.lc, r.md.Attrs)
	s.rc.Add(id, r)
	return r, nil
}

func (s *storage) NewRelation(id string, md metadata.Metadata) error {
	s.Lock()
	defer s.Unlock()
	data, err := encoding.Encode(md)
	if err != nil {
		return err
	}
	bat, err := s.db.NewBatch()
	if err != nil {
		return err
	}
	if !md.IsE {
		if md.Attrs[0].Name != index.SEQ && md.Attrs[0].Type != types.T_uint64 {
			return fmt.Errorf("need attribute seq(uint64)")
		}
		if err := bat.Set([]byte(Item), []byte(id)); err != nil {
			bat.Cancel()
			return err
		}
	} else {
		if md.Attrs[0].Name != index.SEQ && md.Attrs[0].Type != types.T_uint64 {
			return fmt.Errorf("need attribute seq(uint64)")
		}
		if md.Attrs[1].Name != index.ID && md.Attrs[1].Type != types.T_uint64 {
			return fmt.Errorf("need attribute id(uint64)")
		}
		if err := bat.Set([]byte(Event), []byte(id)); err != nil {
			bat.Cancel()
			return err
		}
	}
	if err := bat.Set(metadata.Mkey(id), data); err != nil {
		bat.Cancel()
		return err
	}
	return bat.Commit()
}

func (r *relation) Destroy() error {
	r.Lock()
	defer r.Unlock()
	return nil
}

func (r *relation) IsEvent() bool {
	r.RLock()
	defer r.RUnlock()
	return r.md.IsE
}

func (r *relation) IdMap() (bsi.Bsi, error) {
	r.RLock()
	defer r.RUnlock()
	return r.idx.IdMap()
}

func (r *relation) AddTuples(ts []interface{}) error {
	r.Lock()
	defer r.Unlock()
	return r.idx.AddTuples(ts)
}

func (r *relation) Eq(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	return r.idx.Eq(attr, v)
}

func (r *relation) Ne(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	return r.idx.Ne(attr, v)
}

func (r *relation) Lt(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	return r.idx.Lt(attr, v)
}

func (r *relation) Le(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	return r.idx.Le(attr, v)
}

func (r *relation) Gt(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	return r.idx.Gt(attr, v)
}

func (r *relation) Ge(attr string, v value.Value) (*roaring.Bitmap, error) {
	r.RLock()
	defer r.RUnlock()
	return r.idx.Ge(attr, v)
}
