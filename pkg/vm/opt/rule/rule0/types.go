package rule0

import (
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
)

type rule struct {
	flg bool // whether there is false
	mp  map[int32]bool
	stg storage.Storage
	c   context.Context
}
