package context

import (
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/storage/metadata"
)

func New(cli client.Client, stg storage.Storage) *context {
	return &context{
		cli: cli,
		stg: stg,
	}
}

func (c *context) Client() client.Client {
	return c.cli
}

func (c *context) IsIndex(name, id string) (bool, error) {
	{
		r, err := c.stg.Relation(metadata.Ikey(id))
		switch {
		case err == nil:
			attrs := r.Metadata().Attrs
			for _, attr := range attrs {
				if attr.Name == name {
					return attr.Index, nil
				}
			}
		case err == NotExist:
		default:
			return false, err
		}
	}
	r, err := c.stg.Relation(metadata.Ikey(id))
	if err != nil {
		return false, err
	}
	attrs := r.Metadata().Attrs
	for _, attr := range attrs {
		if attr.Name == name {
			return attr.Index, nil
		}
	}
	return false, NotExist
}

func (c *context) AttributeBelong(name, id string) (string, error) {
	{
		_, err := c.attributeType(name, metadata.Ikey(id))
		switch {
		case err == nil:
			return metadata.Ikey(id), nil
		case err == NotExist:
		default:
			return "", err
		}
	}
	if _, err := c.attributeType(name, metadata.Ekey(id)); err != nil {
		return "", err
	}
	return metadata.Ekey(id), nil
}

func (c *context) AttributeType(name, id string) (uint32, error) {
	{
		typ, err := c.attributeType(name, metadata.Ikey(id))
		switch {
		case err == nil:
			return typ, nil
		case err == NotExist:
		default:
			return 0, err
		}
	}
	return c.attributeType(name, metadata.Ekey(id))
}

func (c *context) attributeType(name, id string) (uint32, error) {
	r, err := c.stg.Relation(id)
	if err != nil {
		return 0, err
	}
	attrs := r.Metadata().Attrs
	for _, attr := range attrs {
		if attr.Name == name {
			return attr.Type, nil
		}
	}
	return 0, NotExist
}
