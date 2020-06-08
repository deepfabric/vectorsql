package opt

import (
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	"github.com/deepfabric/vectorsql/pkg/vm/opt/rule"
)

type optimizer struct {
	rs  []rule.Rule
	c   context.Context
	stg storage.Storage
}
