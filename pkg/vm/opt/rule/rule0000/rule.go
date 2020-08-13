package rule0000

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/extend/overload"
	"github.com/deepfabric/vectorsql/pkg/vm/filter"
	"github.com/deepfabric/vectorsql/pkg/vm/filter/ck"
	Rule "github.com/deepfabric/vectorsql/pkg/vm/opt/rule"
	"github.com/deepfabric/vectorsql/pkg/vm/opt/rule/bm"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

func New(c context.Context, stg storage.Storage) Rule.Rule {
	return &rule{c: c, stg: stg}
}

func (r *rule) Match(e extend.Extend) bool {
	return true
}

func (r *rule) Rewrite(e extend.Extend, id string) (filter.Filter, filter.Filter, error) {
	bs, qs, err := r.disintegration(e, []bm.Bm{}, []string{}, id)
	if err != nil {
		return nil, nil, err
	}
	return ck.New(r.c.Client(), genQuery(bs, qs)), nil, nil
}

func (r *rule) disintegration(e extend.Extend, bs []bm.Bm, qs []string, id string) ([]bm.Bm, []string, error) {
	switch v := e.(type) {
	case *value.Bool:
		query, err := r.genResult(v, id)
		if err != nil {
			return nil, nil, err
		}
		qs = append(qs, query)
		bs = append(bs, bm.Bm{Name: fmt.Sprintf("bm%v", r.cnt)})
		r.cnt++
	case *extend.ParenExtend:
		var err error
		var bt []bm.Bm
		if bt, qs, err = r.disintegration(v.E, []bm.Bm{}, qs, id); err != nil {
			return nil, nil, err
		}
		bs = append(bs, bm.Bm{Bs: bt})
	case *extend.BinaryExtend:
		var err error
		if bs, qs, err = r.disintegrationBinary(v, bs, qs, id); err != nil {
			return nil, nil, err
		}
	}
	return bs, qs, nil
}

func (r *rule) disintegrationBinary(e *extend.BinaryExtend, bs []bm.Bm, qs []string, id string) ([]bm.Bm, []string, error) {
	switch e.Op {
	case overload.EQ, overload.LT, overload.GT, overload.LE,
		overload.GE, overload.NE, overload.Like, overload.NotLike:
		query, err := r.genResult(e, id)
		if err != nil {
			return nil, nil, err
		}
		qs = append(qs, query)
		bs = append(bs, bm.Bm{Name: fmt.Sprintf("bm%v", r.cnt)})
		r.cnt++
		return bs, qs, nil
	case overload.Or:
		var err error
		if bs, qs, err = r.disintegration(e.Left, bs, qs, id); err != nil {
			return nil, nil, err
		}
		bs[len(bs)-1].IsOr = true
		if bs, qs, err = r.disintegration(e.Right, bs, qs, id); err != nil {
			return nil, nil, err
		}
		return bs, qs, nil
	case overload.And:
		var err error
		if bs, qs, err = r.disintegration(e.Left, bs, qs, id); err != nil {
			return nil, nil, err
		}
		if bs, qs, err = r.disintegration(e.Right, bs, qs, id); err != nil {
			return nil, nil, err
		}
		return bs, qs, nil
	}
	return nil, nil, errors.New("extend must be a boolean expression")
}

func (r *rule) genResult(e extend.Extend, id string) (string, error) {
	ts, err := r.extendBelong(e, id)
	if err != nil {
		return "", err
	}
	switch len(ts) {
	case 0:
		return fmt.Sprintf("(SELECT groupBitmapState(uid) FROM people WHERE %s) AS bm%v", e, r.cnt), nil
	case 1:
		return fmt.Sprintf("(SELECT groupBitmapState(uid) FROM %s WHERE %s) AS bm%v", ts[0], e, r.cnt), nil
	}
	return "", fmt.Errorf("'%s' unsupport now", e)
}

func (r *rule) extendBelong(e extend.Extend, id string) ([]string, error) {
	if attrs := e.Attributes(); len(attrs) > 0 {
		mp := make(map[string]struct{})
		for _, attr := range attrs {
			name, err := r.c.AttributeBelong(attr, id)
			if err != nil {
				return nil, err
			}
			if _, ok := mp[name]; !ok {
				mp[name] = struct{}{}
			}
		}
		rs := make([]string, 0, len(mp))
		for k, _ := range mp {
			rs = append(rs, k)
		}
		return rs, nil
	}
	return nil, nil
}

func genQuery(bs []bm.Bm, qs []string) string {
	var buf bytes.Buffer

	for i, q := range qs {
		if i == 0 {
			buf.WriteString(fmt.Sprintf("WITH %s", q))
		} else {
			buf.WriteString(fmt.Sprintf(", %s", q))
		}
	}
	buf.WriteString(fmt.Sprintf(" SELECT CAST(%s AS String) AS result", bm.Gen(bs)))
	return buf.String()
}
