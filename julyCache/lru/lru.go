package lru

import "container/list"

type Cache struct {
	maxBytes, nBytes int64
	ll               *list.List
	cache            map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look-ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if _, ok = c.cache[key]; ok {
		ele := c.cache[key]
		c.ll.MoveToFront(ele)
		value = ele.Value.(*entry).value
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		node := ele.Value.(*entry)
		delete(c.cache, node.key)
		c.nBytes -= int64(len(node.key)) + int64(node.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(node.key, node.value)
		}
	}
}

// Put adds or update a value to the cache.
func (c *Cache) Put(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		node := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(node.value.Len())
		node.value = value
	} else {
		ele = c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.removeOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

// entry the node of list
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}
