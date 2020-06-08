package main

import (
	"fmt"
	"os"

	"github.com/deepfabric/vectorsql/pkg/sql/parser"
)

func main() {
	stmt, err := parser.Parse(os.Args[1])
	fmt.Printf("%v: %v\n", stmt, err)
}
