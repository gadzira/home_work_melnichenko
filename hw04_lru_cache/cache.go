package hw04_lru_cache //nolint:golint,stylecheck

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
	getCap() int
	getLen() int
}

type lruCache struct {
	capacity int
	queue    list
	items    map[Key]*cacheItem
}

type cacheItem struct {
	key   Key
	Value interface{}
}

type Key string

func (l *lruCache) Set(key Key, value interface{}) bool {
	// If element exists, just update the value and move to the front of the queue
	if val, exist := l.items[key]; exist {
		l.items[key] = &cacheItem{Value: value}

		l.queue.MoveToFront(val)
		// l.queue.PushFront(&cacheItem{key: key, Value: &Node{Value: value}})
		return true
	}
	// If the length of the queue more than the capacity of the cache,
	// remove the last element of the queue and the value from the map
	if l.queue.Len() >= l.capacity {
		lastElementOfQueue := l.queue.Back()
		nodeForRemove := lastElementOfQueue.Value.(*cacheItem)
		delete(l.items, nodeForRemove.key)
		l.Clear()
	}
	// If element absent of the cache - add at the map and move to the front of the queue
	l.items[key] = &cacheItem{key: key, Value: value}
	l.queue.PushFront(&cacheItem{key: key, Value: &Node{Value: value}})
	return false
}

// When you get an element from the cache and element exists at the cache - move an element to the front and return true.
// Else return false and nil
func (l *lruCache) Get(key Key) (interface{}, bool) {
	if value, ok := l.items[key]; ok {
		l.queue.PushFront(value)
		// fmt.Println(value)
		return l.items[key].Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	lqb := l.queue.Back()
	l.queue.Remove(lqb)
}

func (l *lruCache) getCap() int {
	return l.capacity
}

func (l *lruCache) getLen() int {
	return l.queue.Lenght
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    list{},
		items:    make(map[Key]*cacheItem),
	}
}
