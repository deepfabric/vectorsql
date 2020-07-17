package cache

import (
	"container/list"
	"sync"
)

type Cache interface {
	Del(string) error
	Get(string) (interface{}, bool)
	Set(string, interface{}, []byte) error
}

type entry struct {
	n int
	k string
	v interface{}
}

type cache struct {
	sync.Mutex
	size  int
	limit int
	lt    *list.List
	mp    map[string]*list.Element
}
