package hw04_lru_cache //nolint:golint,stylecheck

import "fmt"

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    *list
	items    map[Key]*Node
}

type cacheItem struct {
	key   Key
	Value interface{}
}

type Key string

func (l *lruCache) Set(key Key, value interface{}) bool {
	// If element exists, just update the value and move to the front of the queue
	if val, exist := l.items[key]; exist {
		fmt.Println("Some element exixt?", val.Value)
		l.queue.MoveToFront(val)
		val.Value.(*cacheItem).Value = value
		return true
	}
	if l.queue.Len() >= l.capacity {
		lastElementOfQueue := l.queue.back
		nodeForRemove := lastElementOfQueue.Value.(*cacheItem)
		fmt.Println("Delete element of map by key:", nodeForRemove.key)
		delete(l.items, nodeForRemove.key) // remove key from map
		fmt.Println("Remove node from queue:", lastElementOfQueue.Value)
		l.queue.Remove(lastElementOfQueue) // remove node from queue
	}
	// If element absent of the cache - add at the map and move to the front of the queue
	newCacheItem := l.queue.PushFront(&cacheItem{key: key, Value: value})
	fmt.Println("newCacheItem:", newCacheItem.Value)
	l.items[key] = newCacheItem
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if value, ok := l.items[key]; ok {
		l.queue.MoveToFront(value)
		return l.items[key].Value.(*cacheItem).Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = &list{}
	l.items = make(map[Key]*Node, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    &list{},
		items:    make(map[Key]*Node),
	}
}
