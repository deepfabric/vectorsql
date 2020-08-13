package op

import (
	"bytes"
	"fmt"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/vm/bv"
)

func (o *OP) Result(log logger.Log, b bv.BV, cli client.Client, vec []float32) ([][]string, error) {
	var mp *roaring.Bitmap

	if len(vec) != 512 {
		return nil, fmt.Errorf("illegal vector '%v'", vec)
	}
	switch {
	case o.Cf != nil && o.If == nil:
		if mq, err := o.Cf.Bitmap(); err != nil {
			return nil, err
		} else {
			mp = mq
		}
	case o.Cf == nil && o.If != nil:
		if mq, err := o.If.Bitmap(); err != nil {
			return nil, err
		} else {
			mp = mq
		}
	case o.Cf != nil && o.If != nil:
		if mq, err := o.Cf.Bitmap(); err != nil {
			return nil, err
		} else {
			mp = mq
		}
		if mq, err := o.If.Bitmap(); err != nil {
			return nil, err
		} else {
			mp = roaring.FastAnd(mp, mq)
		}
	}
	switch {
	case o.T != nil && o.T.IsF:
		t := time.Now()
		rp, vs, err := b.Fvectors(int64(o.T.Num), vec)
		if err != nil {
			return nil, err
		}
		{
			log.Debugf("vector process: %v\n", time.Now().Sub(t))
		}
		if mp != nil {
			mp.And(rp)
		} else {
			mp = rp
		}
		is := mp.ToArray()
		switch {
		case len(vs) > 0 && len(is) > 0:
			sel := o.N.Relation.(*tree.SelectClause)
			from := sel.From
			sel.From = nil
			sql := fmt.Sprintf("%s FROM", o.N.String())
			sel.From = from
			sql += fmt.Sprintf(" (WITH %v AS xids", slice2String(vs))
			{
				sel.Sel = append(tree.SelectExprs{&tree.SelectExpr{
					As: tree.Name("no"),
					E:  &tree.Index{},
				}}, sel.Sel...)
				o.N.Relation = sel
			}
			sql += fmt.Sprintf(" %s WHERE xid IN xids AND uid IN %s ORDER BY no)", o.N, slice2String32(is))
			return cli.Query(sql)
		case len(vs) == 0 && len(is) > 0:
			return nil, nil
		case len(vs) > 0 && len(is) == 0:
			sel := o.N.Relation.(*tree.SelectClause)
			from := sel.From
			sel.From = nil
			sql := fmt.Sprintf("%s FROM", o.N.String())
			sel.From = from
			sql += fmt.Sprintf(" (WITH %v AS xids", slice2String(vs))
			{
				sel.Sel = append(tree.SelectExprs{&tree.SelectExpr{
					As: tree.Name("no"),
					E:  &tree.Index{},
				}}, sel.Sel...)
				o.N.Relation = sel
			}
			sql += fmt.Sprintf(" %s WHERE xid IN xids ORDER BY no)", o.N)
			return cli.Query(sql)
		}
		return nil, nil
	case o.T != nil && !o.T.IsF:
		t := time.Now()
		_, vs, err := b.Vectors(int64(o.T.Num), mp, vec)
		if err != nil {
			return nil, err
		}
		{
			log.Debugf("vector process: %v\n", time.Now().Sub(t))
		}
		if len(vs) > 0 {
			sel := o.N.Relation.(*tree.SelectClause)
			from := sel.From
			sel.From = nil
			sql := fmt.Sprintf("%s FROM ", o.N.String())
			sel.From = from
			sql += fmt.Sprintf(" (WITH %v AS xids", slice2String(vs))
			{
				sel := o.N.Relation.(*tree.SelectClause)
				sel.Sel = append(tree.SelectExprs{&tree.SelectExpr{
					As: tree.Name("no"),
					E:  &tree.Index{},
				}}, sel.Sel...)
				o.N.Relation = sel
			}
			sql += fmt.Sprintf(" %s WHERE xid IN xids ORDER BY no)", o.N)
			return cli.Query(sql)
		}
		return nil, nil
	default:
		{
			log.Debugf("query: '%v'\n", o.N.String())
		}
		if is := mp.ToArray(); len(is) > 0 {
			return cli.Query(o.N.String() + fmt.Sprintf(" WHERE uid IN %s", slice2String32(is)))
		} else {
			return cli.Query(o.N.String())
		}
	}
}

func slice2String(is []uint64) string {
	var buf bytes.Buffer

	buf.WriteByte('[')
	for i, v := range is {
		if i > 0 {
			buf.WriteString(fmt.Sprintf(", %v", v))
		} else {
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

func slice2String32(is []uint32) string {
	var buf bytes.Buffer

	buf.WriteByte('[')
	for i, v := range is {
		if i > 0 {
			buf.WriteString(fmt.Sprintf(", %v", v))
		} else {
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	}
	buf.WriteByte(']')
	return buf.String()
}
