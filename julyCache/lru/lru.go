package lru

import (
	"container/list"
	"time"
)

type Cache struct {
	maxBytes, nBytes int64
	ll               *list.List
	cache            map[string]*list.Element
	expireCache      map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	cache := &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
	cache.DeleteExpired()
	return cache
}

// Get look-ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if _, ok = c.cache[key]; ok {
		ele := c.cache[key]
		c.ll.MoveToFront(ele)
		entity := ele.Value.(*entry)
		if entity.ddl > 0 && entity.ddl < time.Now().Unix() {
			c.remove(key)
			return nil, false
		}
		value = entity.value
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		node := ele.Value.(*entry)
		delete(c.cache, node.key)
		delete(c.expireCache, node.key)
		c.nBytes -= int64(len(node.key)) + int64(node.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(node.key, node.value)
		}
	}
}

// remove the item of a specific key
func (c *Cache) remove(key string) {
	element := c.cache[key]
	c.ll.MoveToBack(element)
	c.RemoveOldest()
}

// Put adds or update a value to the cache.
func (c *Cache) Put(key string, value Value, expire ...int64) {
	var ddl int64 = -1
	if len(expire) != 0 && expire[0] > 0 {
		ddl = time.Now().Unix() + expire[0]
	}
	ele, ok := c.cache[key]
	if ok {
		c.ll.MoveToFront(ele)
		node := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(node.value.Len())
		node.value = value
		node.ddl = ddl
	} else {
		ele = c.ll.PushFront(&entry{key, value, ddl})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	if len(expire) != 0 && expire[0] > 0 {
		if c.expireCache == nil {
			c.expireCache = make(map[string]*list.Element)
		}
		c.expireCache[key] = ele
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) DeleteExpired() {
	go func() {
		for true {
			if c.expireCache == nil {
				time.Sleep(1 * time.Second)
				continue
			}
			count := 20
			expired := 0
			for _, v := range c.expireCache {
				if count <= 0 {
					break
				}
				e := v.Value.(*entry)
				if e.ddl <= time.Now().Unix() {
					expired++
					c.remove(e.key)
				}
				count--
			}
			if expired < 5 {
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

// entry the node of list
type entry struct {
	key   string
	value Value
	ddl   int64
}

type Value interface {
	Len() int
}
