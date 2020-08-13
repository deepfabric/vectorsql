package rule0

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/extend/overload"
	"github.com/deepfabric/vectorsql/pkg/vm/filter"
	"github.com/deepfabric/vectorsql/pkg/vm/filter/ck"
	"github.com/deepfabric/vectorsql/pkg/vm/filter/index"
	"github.com/deepfabric/vectorsql/pkg/vm/filter/index/ifilter"
	Rule "github.com/deepfabric/vectorsql/pkg/vm/opt/rule"
	"github.com/deepfabric/vectorsql/pkg/vm/opt/rule/bm"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

func New(c context.Context, stg storage.Storage) Rule.Rule {
	return &rule{
		c:   c,
		stg: stg,
		mp: map[uint32]bool{
			types.T_int8:      true,
			types.T_int16:     true,
			types.T_int32:     true,
			types.T_int64:     true,
			types.T_uint8:     true,
			types.T_uint16:    true,
			types.T_uint32:    true,
			types.T_uint64:    true,
			types.T_float32:   true,
			types.T_float64:   true,
			types.T_string:    false,
			types.T_timestamp: true,
		},
		mq: map[uint32]bool{
			types.T_int8:      true,
			types.T_int16:     true,
			types.T_int32:     true,
			types.T_int64:     true,
			types.T_uint8:     true,
			types.T_uint16:    true,
			types.T_uint32:    true,
			types.T_uint64:    true,
			types.T_float32:   true,
			types.T_float64:   true,
			types.T_string:    true,
			types.T_timestamp: true,
		},
	}
}

func (r *rule) Match(e extend.Extend) bool {
	return e.IsAndOnly()
}

func (r *rule) Rewrite(e extend.Extend, id string) (filter.Filter, filter.Filter, error) {
	if r.flg {
		return nil, nil, nil
	}
	mp, mq, err := r.disintegration(e, id)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case mp == nil && mq != nil:
		fl, err := r.genIndexFilter(mq)
		if err != nil {
			return nil, nil, err
		}
		return nil, fl, nil
	case mp != nil && mq == nil:
		return ck.New(r.c.Client(), genQuery(mp)), nil, nil
	case mp != nil && mq != nil:
		fl, err := r.genIndexFilter(mq)
		if err != nil {
			return nil, nil, err
		}
		return ck.New(r.c.Client(), genQuery(mp)), fl, nil
	}
	return nil, nil, nil
}

func (r *rule) disintegration(e extend.Extend, id string) (map[string]extend.Extend, map[string][]*ifilter.Condition, error) {
	if !e.IsLogical() {
		return nil, nil, errors.New("extend must be a boolean expression")
	}
	switch v := e.(type) {
	case *value.Bool:
		if !value.MustBeBool(v) {
			r.flg = true
		}
		return nil, nil, nil
	case *extend.ParenExtend:
		return r.disintegration(v.E, id)
	case *extend.BinaryExtend:
		return r.disintegrationBinary(v, id)
	}
	return nil, nil, errors.New("extend must be a boolean expression")
}

func (r *rule) disintegrationBinary(e *extend.BinaryExtend, id string) (map[string]extend.Extend, map[string][]*ifilter.Condition, error) {
	switch e.Op {
	case overload.EQ:
		c, err := r.buildEQ(e, id)
		if err != nil {
			return nil, nil, err
		}
		return r.genResult(e, c, id)
	case overload.LT:
		c, err := r.buildLT(e, id)
		if err != nil {
			return nil, nil, err
		}
		return r.genResult(e, c, id)
	case overload.GT:
		c, err := r.buildGT(e, id)
		if err != nil {
			return nil, nil, err
		}
		return r.genResult(e, c, id)
	case overload.LE:
		c, err := r.buildLE(e, id)
		if err != nil {
			return nil, nil, err
		}
		return r.genResult(e, c, id)
	case overload.GE:
		c, err := r.buildGE(e, id)
		if err != nil {
			return nil, nil, err
		}
		return r.genResult(e, c, id)
	case overload.NE:
		c, err := r.buildNE(e, id)
		if err != nil {
			return nil, nil, err
		}
		return r.genResult(e, c, id)
	case overload.And:
		lp, lq, err := r.disintegration(e.Left, id)
		if err != nil {
			return nil, nil, err
		}
		rp, rq, err := r.disintegration(e.Right, id)
		if err != nil {
			return nil, nil, err
		}
		switch {
		case lp == nil && rp == nil:
			for k, v := range rq {
				lq[k] = append(lq[k], v...)
			}
			return nil, lq, nil
		case lp != nil && rp != nil:
			for k, rv := range rp {
				if lv, ok := lp[k]; ok {
					lp[k] = &extend.BinaryExtend{
						Left:  lv,
						Right: rv,
						Op:    overload.And,
					}
				} else {
					lp[k] = rv
				}
			}
			return lp, nil, nil
		case lp == nil && rp != nil:
			return rp, lq, nil
		default:
			return lp, rq, nil
		}
	case overload.Like, overload.NotLike:
		ts, err := r.extendBelong(e, id)
		if err != nil {
			return nil, nil, err
		}
		if len(ts) != 1 {
			return nil, nil, fmt.Errorf("unsupport '%s' now", e)
		}
		mp := make(map[string]extend.Extend)
		mp[ts[0]] = e
		return mp, nil, nil
	}
	return nil, nil, errors.New("extend must be a boolean expression")
}

func (r *rule) buildEQ(e *extend.BinaryExtend, id string) (*ifilter.Condition, error) {
	left, right := e.Left, e.Right
	if lv, ok := left.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(lv.Name, id)
		if err != nil {
			return nil, err
		}
		if rv, ok := right.(value.Value); ok && r.mq[typ] && r.typeCheck(typ, rv.ResolvedType(), lv.Name, id) {
			return &ifilter.Condition{Op: ifilter.EQ, Name: lv.Name, Val: typeCast(typ, rv)}, nil
		}
	}
	if rv, ok := right.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(rv.Name, id)
		if err != nil {
			return nil, err
		}
		if lv, ok := left.(value.Value); ok && r.mq[typ] && r.typeCheck(typ, lv.ResolvedType(), rv.Name, id) {
			return &ifilter.Condition{Op: ifilter.EQ, Name: rv.Name, Val: typeCast(typ, lv)}, nil
		}
	}
	return nil, nil
}

func (r *rule) buildNE(e *extend.BinaryExtend, id string) (*ifilter.Condition, error) {
	left, right := e.Left, e.Right
	if lv, ok := left.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(lv.Name, id)
		if err != nil {
			return nil, err
		}
		if rv, ok := right.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, rv.ResolvedType(), lv.Name, id) {
			return &ifilter.Condition{Op: ifilter.NE, Name: lv.Name, Val: typeCast(typ, rv)}, nil
		}
	}
	if rv, ok := right.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(rv.Name, id)
		if err != nil {
			return nil, err
		}
		if lv, ok := left.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, lv.ResolvedType(), rv.Name, id) {
			return &ifilter.Condition{Op: ifilter.NE, Name: rv.Name, Val: typeCast(typ, lv)}, nil
		}
	}
	return nil, nil
}

func (r *rule) buildLT(e *extend.BinaryExtend, id string) (*ifilter.Condition, error) {
	left, right := e.Left, e.Right
	if lv, ok := left.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(lv.Name, id)
		if err != nil {
			return nil, err
		}
		if rv, ok := right.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, rv.ResolvedType(), lv.Name, id) {
			return &ifilter.Condition{Op: ifilter.LT, Name: lv.Name, Val: typeCast(typ, rv)}, nil
		}
	}
	if rv, ok := right.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(rv.Name, id)
		if err != nil {
			return nil, err
		}
		if lv, ok := left.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, lv.ResolvedType(), rv.Name, id) {
			return &ifilter.Condition{Op: ifilter.GE, Name: rv.Name, Val: typeCast(typ, lv)}, nil
		}
	}
	return nil, nil
}

func (r *rule) buildLE(e *extend.BinaryExtend, id string) (*ifilter.Condition, error) {
	left, right := e.Left, e.Right
	if lv, ok := left.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(lv.Name, id)
		if err != nil {
			return nil, err
		}
		if rv, ok := right.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, rv.ResolvedType(), lv.Name, id) {
			return &ifilter.Condition{Op: ifilter.LE, Name: lv.Name, Val: typeCast(typ, rv)}, nil
		}
	}
	if rv, ok := right.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(rv.Name, id)
		if err != nil {
			return nil, err
		}
		if lv, ok := left.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, lv.ResolvedType(), rv.Name, id) {
			return &ifilter.Condition{Op: ifilter.GT, Name: rv.Name, Val: typeCast(typ, lv)}, nil
		}
	}
	return nil, nil
}

func (r *rule) buildGT(e *extend.BinaryExtend, id string) (*ifilter.Condition, error) {
	left, right := e.Left, e.Right
	if lv, ok := left.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(lv.Name, id)
		if err != nil {
			return nil, err
		}
		if rv, ok := right.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, rv.ResolvedType(), lv.Name, id) {
			return &ifilter.Condition{Op: ifilter.GT, Name: lv.Name, Val: typeCast(typ, rv)}, nil
		}
	}
	if rv, ok := right.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(rv.Name, id)
		if err != nil {
			return nil, err
		}
		if lv, ok := left.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, lv.ResolvedType(), rv.Name, id) {
			return &ifilter.Condition{Op: ifilter.LE, Name: rv.Name, Val: typeCast(typ, lv)}, nil
		}
	}
	return nil, nil
}

func (r *rule) buildGE(e *extend.BinaryExtend, id string) (*ifilter.Condition, error) {
	left, right := e.Left, e.Right
	if lv, ok := left.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(lv.Name, id)
		if err != nil {
			return nil, err
		}
		if rv, ok := right.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, rv.ResolvedType(), lv.Name, id) {
			return &ifilter.Condition{Op: ifilter.GE, Name: lv.Name, Val: typeCast(typ, rv)}, nil
		}
	}
	if rv, ok := right.(*extend.Attribute); ok {
		typ, err := r.c.AttributeType(rv.Name, id)
		if err != nil {
			return nil, err
		}
		if lv, ok := left.(value.Value); ok && r.mp[typ] && r.typeCheck(typ, lv.ResolvedType(), rv.Name, id) {
			return &ifilter.Condition{Op: ifilter.LT, Name: rv.Name, Val: typeCast(typ, lv)}, nil
		}
	}
	return nil, nil
}

func (r *rule) genIndexFilter(mp map[string][]*ifilter.Condition) (filter.Filter, error) {
	fs := make([]ifilter.Filter, 0, len(mp))
	for k, v := range mp {
		r, err := r.stg.Relation(k)
		if err != nil {
			return nil, err
		}
		fs = append(fs, ifilter.New(v, r))
	}
	return index.New(fs), nil
}

func (r *rule) genResult(e extend.Extend, c *ifilter.Condition, id string) (map[string]extend.Extend, map[string][]*ifilter.Condition, error) {
	ts, err := r.extendBelong(e, id)
	if err != nil {
		return nil, nil, err
	}
	if c != nil {
		mp := make(map[string][]*ifilter.Condition)
		mp[ts[0]] = []*ifilter.Condition{c}
		return nil, mp, nil
	}
	if len(ts) != 1 {
		return nil, nil, fmt.Errorf("'%s' unsupport now", e)
	}
	mp := make(map[string]extend.Extend)
	mp[ts[0]] = e
	return mp, nil, nil
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

func typeCast(x uint32, v value.Value) value.Value {
	switch {
	case x == types.T_int8 && v.ResolvedType() == types.T_int:
		return value.NewInt8(int8(value.MustBeInt(v) & 0xFF))
	case x == types.T_int16 && v.ResolvedType() == types.T_int:
		return value.NewInt16(int16(value.MustBeInt(v) & 0xFFFF))
	case x == types.T_int32 && v.ResolvedType() == types.T_int:
		return value.NewInt32(int32(value.MustBeInt(v) & 0xFFFFFFFF))
	case x == types.T_int64 && v.ResolvedType() == types.T_int:
		return value.NewInt64(value.MustBeInt(v))
	case x == types.T_uint8 && v.ResolvedType() == types.T_int:
		return value.NewUint8(uint8(value.MustBeInt(v) & 0xFF))
	case x == types.T_uint16 && v.ResolvedType() == types.T_int:
		return value.NewUint16(uint16(value.MustBeInt(v) & 0xFFFF))
	case x == types.T_uint32 && v.ResolvedType() == types.T_int:
		return value.NewUint32(uint32(value.MustBeInt(v) & 0xFFFFFFFF))
	case x == types.T_uint64 && v.ResolvedType() == types.T_int:
		return value.NewUint64(uint64(value.MustBeInt(v)))
	case x == types.T_float32 && v.ResolvedType() == types.T_float:
		return value.NewFloat32(float32(value.MustBeFloat(v)))
	case x == types.T_float64 && v.ResolvedType() == types.T_float:
		return value.NewFloat64(value.MustBeFloat(v))
	case x == types.T_timestamp && v.ResolvedType() == types.T_int:
		return value.NewTimestamp(time.Unix(value.MustBeInt(v), 0))
	case x == types.T_timestamp && v.ResolvedType() == types.T_string:
		rv, _ := value.ParseTimestamp(value.MustBeString(v))
		return rv
	}
	return v
}

func (r *rule) typeCheck(x uint32, typ types.T, name, id string) bool {
	y := uint32(typ)
	switch x {
	case types.T_int8:
		return y == x || y == types.T_int
	case types.T_int16:
		return y == x || y == types.T_int
	case types.T_int32:
		return y == x || y == types.T_int
	case types.T_int64:
		return y == x || y == types.T_int
	case types.T_uint8:
		return y == x || y == types.T_int
	case types.T_uint16:
		return y == x || y == types.T_int
	case types.T_uint32:
		return y == x || y == types.T_int
	case types.T_uint64:
		return y == x || y == types.T_int
	case types.T_float32:
		return y == x || y == types.T_float
	case types.T_float64:
		return y == x || y == types.T_float
	case types.T_timestamp:
		return y == x || y == types.T_int || y == types.T_string
	case types.T_string:
		ok, _ := r.c.IsIndex(name, id)
		return ok && y == types.T_string
	}
	return x == y
}

func genQuery(mp map[string]extend.Extend) string {
	var buf bytes.Buffer

	cnt := 0
	bs := make([]bm.Bm, 0, len(mp))
	for k, v := range mp {
		if cnt == 0 {
			buf.WriteString(fmt.Sprintf("WITH (SELECT groupBitmapState(uid) FROM %s where %s) AS bm%v", k, v, cnt))
		} else {
			buf.WriteString(fmt.Sprintf(", (SELECT groupBitmapState(uid) FROM %s where %s) AS bm%v", k, v, cnt))
		}
		bs = append(bs, bm.Bm{Name: fmt.Sprintf("bm%v", cnt)})
		cnt++
	}
	buf.WriteString(fmt.Sprintf(" SELECT CAST(%s AS String) AS result", bm.Gen(bs)))
	return buf.String()
}
