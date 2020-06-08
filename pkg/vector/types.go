package vector

import "github.com/deepfabric/vectorsql/pkg/request"

type Vector interface {
	GetVector(map[string]*request.Part) ([]float32, error)
}

type vector struct {
	url string
}
