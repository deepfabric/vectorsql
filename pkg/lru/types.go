package lru

import (
	"container/list"
)

type LRU interface {
	Del(string) error
	Add(string, interface{}) error
	Get(string) (interface{}, bool)
}

type entry struct {
	k string
	v interface{}
}

type lru struct {
	cnt   int
	limit int
	lt    *list.List
	mp    map[string]*list.Element
}
