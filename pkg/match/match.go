package match

import (
	"regexp"

	"github.com/deepfabric/vectorsql/pkg/lru"
	"github.com/deepfabric/vectorsql/pkg/match/like"
)

func New(lc lru.LRU) *match {
	return &match{lc}
}

func (m *match) Compile(expr string, islike bool) (Regexp, error) {
	var err error
	var reg Regexp

	if v, ok := m.lc.Get(expr); ok {
		return v.(Regexp), nil
	}
	switch {
	case islike:
		if reg, err = like.Compile(expr); err != nil {
			return nil, err
		}
	default:
		if reg, err = regexp.Compile(expr); err != nil {
			return nil, err
		}
	}
	m.lc.Add(expr, reg)
	return reg, nil
}
