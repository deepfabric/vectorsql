package vector

type Vector interface {
	GetVector(map[string][]byte) ([]float32, error)
}

type vector struct {
	url string
}
