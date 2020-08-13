package client

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"log"
	"reflect"
	"unsafe"

	"github.com/RoaringBitmap/roaring"

	_ "github.com/ClickHouse/clickhouse-go"
)

func New(dsn string) (*client, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(20)
	return &client{db}, nil
}

func (c *client) Close() error {
	return c.db.Close()
}

func (c *client) Query(query string) ([][]string, error) {
	var rs [][]string

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	attrs, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(attrs))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			log.Fatal(err)
		}
		var v string
		r := make([]string, len(attrs))
		for i, col := range values {
			if col == nil {
				v = "NULL"
			} else {
				v = string(col)
			}
			r[i] = v
		}
		rs = append(rs, r)
	}
	return rs, nil
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

func (c *client) Bitmap(query string) (*roaring.Bitmap, error) {
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
	return unmarshalBitMap([]byte(r))
}

func unmarshalBitMap(data []byte) (*roaring.Bitmap, error) {
	switch data[0] {
	case 0:
		data = data[1:]
		_, n := binary.Uvarint(data)
		if n < 0 {
			return nil, errors.New("overflow")
		}
		return roaring.BitmapOf(decodeVector(data[n:])...), nil
	case 1:
		data = data[1:]
		_, n := binary.Uvarint(data)
		if n < 0 {
			return nil, errors.New("overflow")
		}
		data = data[n:]
		mp := roaring.New()
		if err := mp.UnmarshalBinary(data); err != nil {
			return nil, err
		}
		return mp, nil
	}
	return nil, nil
}

func decodeVector(v []byte) []uint32 {
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	hp.Len /= 4
	hp.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&hp))
}
