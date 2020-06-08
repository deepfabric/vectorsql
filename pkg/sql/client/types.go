package client

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Client interface {
	Close() error
	Query(string, string) (interface{}, error)
}

type client struct {
	db *sqlx.DB
}

type People struct {
	Seq    uint32    `db:"seq"`
	Gender uint8     `db:"gender"`
	Vec    []float32 `db:"vec"`
	//	Vecs   []reflect.Value `db:"vecs"`
}

type PeopleEvent struct {
	Dt   time.Time `db:"dt"`
	Area string    `db:"area"`
	Seq  uint32    `db:"seq"`
}

type Bitmap struct {
	Result string `db:"result"`
}
