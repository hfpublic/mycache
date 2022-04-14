package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int64
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	if ele, has := c.cache[key]; has {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) RemoveOldset() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		e := ele.Value.(*entry)
		delete(c.cache, e.key)
		c.nbytes -= int64(len(e.key)) + e.value.Len()
		if c.OnEvicted != nil {
			c.OnEvicted(e.key, e.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, has := c.cache[key]; has {
		c.ll.MoveToFront(ele)
		e := ele.Value.(*entry)
		c.nbytes += e.value.Len() - value.Len()
		e.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.nbytes += int64(len(key)) + value.Len()
		c.cache[key] = ele
	}
	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldset()
	}
}
