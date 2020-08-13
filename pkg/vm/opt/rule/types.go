package rule

import (
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/filter"
)

type Rule interface {
	Match(extend.Extend) bool
	Rewrite(extend.Extend, string) (filter.Filter, filter.Filter, error)
}
