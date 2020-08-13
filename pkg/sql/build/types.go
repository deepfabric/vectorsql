package build

import (
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
)

type build struct {
	sql string
	c   context.Context
	stg storage.Storage
}
