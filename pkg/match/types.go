package match

import (
	"regexp"

	"github.com/deepfabric/vectorsql/pkg/lru"
)

type Regexp interface {
	MatchString(string) bool
}

type Match interface {
	Compile(string, bool) (Regexp, error)
}

type match struct {
	lc lru.LRU
}

type regular struct {
	reg *regexp.Regexp
	fn  func(string) bool
}
