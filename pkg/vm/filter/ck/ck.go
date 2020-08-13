package ck

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
)

func New(cli client.Client, query string) *ck {
	return &ck{cli: cli, query: query}
}

func (c *ck) String() string {
	return c.query
}

func (c *ck) Bitmap() (*roaring.Bitmap, error) {
	return c.cli.Bitmap(c.query)
}
