package opt

import (
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/filter"
	"github.com/deepfabric/vectorsql/pkg/vm/opt/rule"
	"github.com/deepfabric/vectorsql/pkg/vm/opt/rule/rule0"
	"github.com/deepfabric/vectorsql/pkg/vm/opt/rule/rule0000"
)

func New(c context.Context, stg storage.Storage) *optimizer {
	rs := make([]rule.Rule, 0, len(Rules))
	for _, r := range Rules {
		rs = append(rs, r(c, stg))
	}
	return &optimizer{c: c, rs: rs, stg: stg}
}

func (o *optimizer) Optimize(e extend.Extend, id string) (filter.Filter, filter.Filter, error) {
	for _, r := range o.rs {
		if r.Match(e) {
			return r.Rewrite(e, id)
		}
	}
	return nil, nil, nil
}

var Rules = []func(context.Context, storage.Storage) rule.Rule{
	rule0.New,
	rule0000.New, // default rule
}
