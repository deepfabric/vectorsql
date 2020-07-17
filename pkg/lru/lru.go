package lru

import "container/list"

func New(limit int) *lru {
	return &lru{
		limit: limit,
		lt:    list.New(),
		mp:    make(map[string]*list.Element),
	}
}

func (c *lru) Del(k string) error {
	if e, ok := c.mp[k]; ok {
		c.cnt--
		c.lt.Remove(e)
		delete(c.mp, k)
	}
	return nil
}

func (c *lru) Add(k string, v interface{}) error {
	if e, ok := c.mp[k]; ok {
		c.lt.MoveToFront(e)
		{
			et := e.Value.(*entry)
			et.v = v
		}
		return nil
	}
	c.mp[k] = c.lt.PushFront(&entry{k, v})
	if c.cnt += 1; c.cnt >= c.limit {
		c.reduce()
	}
	return nil
}

func (c *lru) Get(k string) (interface{}, bool) {
	if e, ok := c.mp[k]; ok {
		c.lt.MoveToFront(e)
		return e.Value.(*entry).v, true
	}
	return nil, false
}

func (c *lru) reduce() {
	for e := c.lt.Back(); e != nil; e = c.lt.Back() {
		if c.cnt < c.limit {
			return
		}
		c.cnt--
		delete(c.mp, e.Value.(*entry).k)
		c.lt.Remove(e)
	}
}
