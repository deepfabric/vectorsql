package cache

import "container/list"

func New(limit int) *cache {
	return &cache{
		limit: limit,
		lt:    list.New(),
		mp:    make(map[string]*list.Element),
	}
}

func (c *cache) Del(k string) error {
	c.Lock()
	err := c.del(k)
	c.Unlock()
	return err
}

func (c *cache) Set(k string, v interface{}, data []byte) error {
	c.Lock()
	err := c.set(k, v, data)
	c.Unlock()
	return err
}

func (c *cache) Get(k string) (interface{}, bool) {
	c.Lock()
	v, ok := c.get(k)
	c.Unlock()
	return v, ok
}

func (c *cache) del(k string) error {
	if e, ok := c.mp[k]; ok {
		c.size -= e.Value.(*entry).n
		c.lt.Remove(e)
		delete(c.mp, k)
	}
	return nil
}

func (c *cache) set(k string, v interface{}, data []byte) error {
	if e, ok := c.mp[k]; ok {
		c.lt.MoveToFront(e)
		{
			et := e.Value.(*entry)
			et.v = v
			c.size += len(data) - et.n
			et.n = len(data)
		}
		return nil
	}
	c.mp[k] = c.lt.PushFront(&entry{len(data), k, v})
	if c.size += len(data); c.size >= c.limit {
		c.reduce()
	}
	return nil
}

func (c *cache) get(k string) (interface{}, bool) {
	if e, ok := c.mp[k]; ok {
		c.lt.MoveToFront(e)
		return e.Value.(*entry).v, true
	}
	return nil, false
}

func (c *cache) reduce() {
	for e := c.lt.Back(); e != nil; e = c.lt.Back() {
		if c.size < c.limit {
			return
		}
		c.size -= e.Value.(*entry).n
		delete(c.mp, e.Value.(*entry).k)
		c.lt.Remove(e)
	}
}
