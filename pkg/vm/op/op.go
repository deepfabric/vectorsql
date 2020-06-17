package op

import (
	"bytes"
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/vm/bv"
)

func (o *OP) Result(log logger.Log, mcpu int, b bv.BV, cli client.Client, vec []float32) (interface{}, error) {
	var mp *roaring.Bitmap

	switch {
	case o.Cf != nil && o.If == nil:
		if mq, err := o.Cf.Bitmap(mcpu); err != nil {
			return nil, err
		} else {
			mp = mq
		}
	case o.Cf == nil && o.If != nil:
		if mq, err := o.If.Bitmap(mcpu); err != nil {
			return nil, err
		} else {
			mp = mq
		}
	case o.Cf != nil && o.If != nil:
		if mq, err := o.Cf.Bitmap(mcpu); err != nil {
			return nil, err
		} else {
			mp = mq
		}
		if mq, err := o.Cf.Bitmap(mcpu); err != nil {
			return nil, err
		} else {
			mp = roaring.FastAnd(mp, mq)
		}
	}
	switch {
	case o.T != nil && o.T.IsF:
		rp, err := b.Fvectors(o.T.Num, vec)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			mp.And(rp)
		} else {
			mp = rp
		}
		if is := mp.ToArray(); len(is) > 0 {
			{
				log.Debugf("query: '%v'\n", o.N.String()+fmt.Sprintf(" WHERE seq IN %s", slice2String(is)))
			}
			return cli.Query(o.N.String()+fmt.Sprintf(" WHERE seq IN %s", slice2String(is)), "item")
		}
		return nil, nil
	case o.T != nil && !o.T.IsF:
		rp, err := b.Vectors(mp, o.T.Num, vec)
		if err != nil {
			return nil, err
		}
		if mp != nil {
			mp.And(rp)
		} else {
			mp = rp
		}
		if is := mp.ToArray(); len(is) > 0 {
			{
				log.Debugf("query: '%v'\n", o.N.String()+fmt.Sprintf(" WHERE seq IN %s", slice2String(is)))
			}
			return cli.Query(o.N.String()+fmt.Sprintf(" WHERE seq IN %s", slice2String(is)), "item")
		}
		return nil, nil
	default:
		{
			log.Debugf("query: '%v'\n", o.N.String())
		}
		return cli.Query(o.N.String(), "item")
	}
}

func slice2String(is []uint32) string {
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
