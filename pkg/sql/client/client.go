package client

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go"
)

func New(dsn string) (*client, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	return &client{db}, nil
}

func (c *client) Close() error {
	return c.db.Close()
}

func (c *client) Exec(query string, args [][]interface{}) error {
	if len(args) == 0 {
		if _, err := c.db.Exec(query); err != nil {
			return err
		}
		return nil
	}
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, _ := tx.Prepare(query)
	defer stmt.Close()
	for _, arg := range args {
		if _, err := stmt.Exec(arg...); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *client) Query(query string, typ string) ([]string, error) {
	switch typ {
	case "item":
		rows, err := c.db.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		attrs, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		var rs []string
		values := make([]sql.RawBytes, len(attrs))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		for rows.Next() {
			if err = rows.Scan(scanArgs...); err != nil {
				log.Fatal(err)
			}
			var r, v string
			for i, col := range values {
				if col == nil {
					v = "NULL"
				} else {
					v = string(col)
				}
				if i == 0 {
					r += fmt.Sprintf("%v", v)
				} else {
					r += fmt.Sprintf(", %v", v)
				}
			}
			rs = append(rs, r)
		}
		return rs, nil
	case "bitmap":
		var r string

		rows, err := c.db.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&r); err != nil {
				return nil, err
			}
			break
		}
		return []string{r}, nil
	default:
		return nil, fmt.Errorf("unknown type '%s'", typ)
	}
}
