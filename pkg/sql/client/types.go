package client

import "database/sql"

type Client interface {
	Close() error
	Exec(string, [][]interface{}) error
	Query(string, string) ([]string, error)
}

type client struct {
	db *sql.DB
}
