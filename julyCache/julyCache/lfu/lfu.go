package lfu

import "container/list"

type Cache struct {
	freqLl map[int]*list.List
	cache  map[string]*list.Element
	// cap is the max size of Cache
	cap     int
	minFreq int
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// New is the Constructor of Cache
func New(capacity int, onEvicted func(key string, value Value)) *Cache {
	lfu := &Cache{
		freqLl:    make(map[int]*list.List),
		cache:     make(map[string]*list.Element),
		cap:       capacity,
		minFreq:   1,
		OnEvicted: onEvicted,
	}
	lfu.freqLl[lfu.minFreq] = list.New()
	return lfu
}

// Get look-ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if _, ok = c.cache[key]; ok {
		ele := c.cache[key]
		value = ele.Value.(*entry).value
		c.nodeUpdate(ele)
	}
	return
}

// Put adds or update a value to the cache.
func (c *Cache) Put(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		ele.Value.(*entry).value = value
		c.nodeUpdate(ele)
	} else {
		if c.Len() >= c.cap {
			c.removeLeast()
		}
		kv := &entry{
			key:   key,
			value: value,
			freq:  1,
		}
		node := c.freqLl[kv.freq].PushFront(kv)
		c.cache[key] = node
		c.minFreq = 1
	}
}

// RemoveOldest removes the least item
func (c *Cache) removeLeast() {
	l := c.freqLl[c.minFreq]
	ele := l.Back()
	if ele != nil {
		l.Remove(ele)
		node := ele.Value.(*entry)
		delete(c.cache, node.key)
		if c.OnEvicted != nil {
			c.OnEvicted(node.key, node.value)
		}
	}
}

//nodeUpdate update node's frequence
func (c *Cache) nodeUpdate(node *list.Element) {
	kv := node.Value.(*entry)
	oldList := c.freqLl[kv.freq]
	oldList.Remove(node)

	if oldList.Len() == 0 && c.minFreq == kv.freq {
		c.minFreq++
	}
	// update frequence
	kv.freq++
	if _, ok := c.freqLl[kv.freq]; !ok {
		c.freqLl[kv.freq] = list.New()
	}
	newList := c.freqLl[kv.freq]
	node = newList.PushFront(kv)
	c.cache[kv.key] = node
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return len(c.cache)
}

// entry the node of a list
type entry struct {
	key   string
	value Value
	freq  int
}

type Value interface {
	Len() int
}
