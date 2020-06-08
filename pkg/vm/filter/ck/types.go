package ck

import (
	"github.com/deepfabric/vectorsql/pkg/sql/client"
)

type ck struct {
	query string
	cli   client.Client
}
