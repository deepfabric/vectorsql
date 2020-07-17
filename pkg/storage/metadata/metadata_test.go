package metadata

import (
	"fmt"
	"log"
	"testing"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/util/encoding"
)

func TestMetadata(t *testing.T) {
	var tm Metadata
	var as []Attribute

	{
		as = append(as, Attribute{types.T_uint8, "age"})
		as = append(as, Attribute{types.T_string, "name"})
	}
	md := Metadata{true, as}
	data, err := encoding.Encode(md)
	if err != nil {
		log.Fatal(err)
	}
	{
		fmt.Printf("md: %v\n", md)
	}
	if err := encoding.Decode(data, &tm); err != nil {
		log.Fatal(err)
	}
	{
		fmt.Printf("tm: %v\n", tm)
	}
}
