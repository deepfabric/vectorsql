package rewrite

import "github.com/deepfabric/vectorsql/pkg/vm/extend"

type Rewrite interface {
	Rewrite(extend.Extend) extend.Extend
}
