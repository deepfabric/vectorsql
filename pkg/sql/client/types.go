package client

import "database/sql"

type Client interface {
	Close() error
	Query(string, string) ([]string, error)
	Insert(string, string, [][]interface{}) error
}

type client struct {
	db *sql.DB
}
