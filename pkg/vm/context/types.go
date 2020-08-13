package context

import (
	"errors"

	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
)

var (
	NotExist = errors.New("not exist")
)

type Context interface {
	Client() client.Client
	IsIndex(string, string) (bool, error)
	AttributeType(string, string) (uint32, error)
	AttributeBelong(string, string) (string, error)
}

type context struct {
	cli client.Client
	stg storage.Storage
}
