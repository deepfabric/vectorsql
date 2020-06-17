package context

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

var Types = map[string]int32{
	"seq":  types.T_uint32,
	"age":  types.T_uint8,
	"sex":  types.T_uint8,
	"area": types.T_string,
}

var Tables = map[string]string{
	"seq":  "people",
	"age":  "people",
	"sex":  "people",
	"area": "people",
}

func New(cli client.Client) *context {
	return &context{
		cli: cli,
		mp:  Types,
		mq:  Tables,
	}
}

func (c *context) Client() client.Client {
	return c.cli
}

func (c *context) AttributeBelong(attr string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if v, ok := c.mq[attr]; ok {
		return v, nil
	}
	return "", fmt.Errorf("attribute '%s' not exist", attr)
}

func (c *context) AttributeType(attr string) (int32, error) {
	c.RLock()
	defer c.RUnlock()
	if v, ok := c.mp[attr]; ok {
		return v, nil
	}
	return -1, fmt.Errorf("attribute '%s' not exist", attr)
}
