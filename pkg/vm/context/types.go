package context

import (
	"sync"

	"github.com/deepfabric/vectorsql/pkg/sql/client"
)

type Context interface {
	Client() client.Client
	AttributeType(string) (int32, error)
	AttributeBelong(string) (string, error)
}

type context struct {
	sync.RWMutex
	cli client.Client
	mp  map[string]int32
	mq  map[string]string
}
