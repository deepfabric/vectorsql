package metadata

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

const (
	mprefix = "_M."    // metadata
	isuffix = "_item"  // item
	esuffix = "_event" // event
)

func init() {
	gob.Register(Metadata{})
	gob.Register(Attribute{})
}

func (a Attribute) String() string {
	return fmt.Sprintf("%s(%s)", a.Name, types.T(a.Type))
}

func Ikey(id string) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteString(isuffix)
	return buf.String()
}

func Ekey(id string) string {
	var buf bytes.Buffer

	buf.WriteString(id)
	buf.WriteString(esuffix)
	return buf.String()
}

func Mkey(id string) []byte {
	var buf bytes.Buffer

	buf.WriteString(mprefix)
	buf.WriteString(id)
	return buf.Bytes()
}
