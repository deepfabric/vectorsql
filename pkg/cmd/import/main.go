package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deepfabric/vectorsql/pkg/sql/client"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: import file...\n")
		return
	}
	cli, err := client.New("tcp://172.19.0.17:9000?username=cdp_user&password=infinivision2019")
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}
	id := loadId()
	for i, j := 1, len(os.Args); i < j; i++ {
		args, ts := loadCsv(id, readFile(os.Args[i]), os.Args[1])
		if err := inject(ts); err != nil {
			log.Fatalf("Failed to inject index: %v\n", err)
		}
		{
			query := "insert into people (seq, sex, age, area) VALUES"
			for _, _ = range args {
				query += " (?, ?, ?. ?)"
			}
			if err := cli.Insert(query, "item", args); err != nil {
				log.Fatalf("Failed to insert to db: %v\n", err)
			}
		}
		id += uint32(len(ts))
		if err := storeId(id); err != nil {
			log.Fatalf("Failed to update id: %v\n", err)
		}
		if err := db.Sync(); err != nil {
			log.Fatalf("Failed to sync index: %v\n", err)
		}
	}
}
