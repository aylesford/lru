package lru

import (
	"container/list"
	"sync"
)

const (
	// DefaultLRUCacheSize defines the default size of cache
	DefaultLRUCacheSize = 1024
)

// option is a model function to provide optional params
type option func(*LRUCache)

// WithCacheSize is used for setting lru cache size
func WithCacheSize(size int) option {
	return func(c *LRUCache) {
		c.size = size
	}
}

// WithExpelCallback is used for setting a function which is called while item is expeled
func WithExpelCallback(f func(key, val interface{})) option {
	return func(c *LRUCache) {
		c.expelCallback = f
	}
}

// LRUCache is a thread-safe cache based on lru
type LRUCache struct {
	size          int
	itemList      *list.List
	itemCache     map[interface{}]*list.Element
	expelCallback func(key, val interface{})
	sync.RWMutex
}

// entry is used for store items of list
type entry struct {
	key interface{}
	val interface{}
}

// NewLRU returns a new LRUCache as you want
func NewLRU(opts ...option) *LRUCache {
	c := &LRUCache{
		size:      DefaultLRUCacheSize,
		itemList:  list.New(),
		itemCache: make(map[interface{}]*list.Element),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Add adds a item into cache
func (c *LRUCache) Add(key, val interface{}) {
	c.Lock()
	defer c.Unlock()

	// Just make it lastest and update value if it's existed
	if item, ok := c.itemCache[key]; ok {
		c.itemList.MoveToFront(item)
		item.Value.(*entry).val = val
		return
	}

	// Add it into cache and make it lastest if it's not existed
	item := c.itemList.PushFront(&entry{key: key, val: val})
	c.itemCache[key] = item

	// Expel oldest item if over size
	if c.itemList.Len() > c.size {
		expelItem := c.itemList.Back()
		c.itemList.Remove(expelItem)
		delete(c.itemCache, expelItem.Value.(*entry).key)

		if c.expelCallback != nil {
			c.expelCallback(expelItem.Value.(*entry).key, expelItem.Value.(*entry).val)
		}
	}
}

// Get looks up a key's value from the cache
func (c *LRUCache) Get(key interface{}) (val interface{}, hit bool) {
	c.RLock()
	item, ok := c.itemCache[key]
	c.RUnlock()

	if ok {
		c.Lock()
		c.itemList.MoveToFront(item)
		c.Unlock()
		return item.Value.(*entry).val, true
	}

	return nil, false
}

// Del removes the provided key from the cache
func (c *LRUCache) Del(key interface{}) {
	c.Lock()
	defer c.Unlock()

	item, ok := c.itemCache[key]
	if ok {
		c.itemList.Remove(item)
		delete(c.itemCache, item.Value.(*entry).key)
	}
}
