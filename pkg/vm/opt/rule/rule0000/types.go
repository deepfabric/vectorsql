package rule0000

import (
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
)

type rule struct {
	cnt int
	c   context.Context
	stg storage.Storage
}
