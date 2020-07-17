package metadata

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

const (
	mprefix = "_M." // metadata
)

func init() {
	gob.Register(Metadata{})
	gob.Register(Attribute{})
}

func (a Attribute) String() string {
	return fmt.Sprintf("%s(%s)", a.Name, types.T(a.Type))
}

func Mkey(id string) []byte {
	var buf bytes.Buffer

	buf.WriteString(mprefix)
	buf.WriteString(id)
	return buf.Bytes()
}
