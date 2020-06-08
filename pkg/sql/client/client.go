package client

import (
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

func New(dsn string) (*client, error) {
	db, err := sqlx.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	return &client{db}, nil
}

func (c *client) Close() error {
	return c.db.Close()
}

func (c *client) Query(query string, table string) (interface{}, error) {
	switch table {
	case "people":
		var rows []People

		if err := c.db.Select(&rows, query); err != nil {
			return nil, err
		}
		return rows, nil
	case "people_events":
		var rows []PeopleEvent

		if err := c.db.Select(&rows, query); err != nil {
			return nil, err
		}
		return rows, nil
	case "bitmap":
		var rows []Bitmap

		if err := c.db.Select(&rows, query); err != nil {
			return nil, err
		}
		return rows, nil
	default:
		return nil, fmt.Errorf("unknown table '%s'", table)
	}
}
