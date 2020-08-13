package client

import (
	"database/sql"

	"github.com/RoaringBitmap/roaring"
)

type Client interface {
	Close() error
	Query(string) ([][]string, error)
	Exec(string, [][]interface{}) error
	Bitmap(string) (*roaring.Bitmap, error)
}

type client struct {
	db *sql.DB
}
