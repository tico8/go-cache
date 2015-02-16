package cache

import (
	"sync"
	"sync/atomic"
	"time"
	"sort"
	"reflect"
	"errors"
)

const (
	// DefaultExpiration expiration of cache item
	DefaultExpiration time.Duration = 60 * 60 * time.Second; // ns 1h
)

// Cache class
type Cache struct {
	*cache
}

// New create instance only once.
// param opt - option
// return arg1 - instance of Cache
func New(opt Option) *Cache {
	c := &cache{
			items:	map[string]*Item{},
			option: &opt,
	}
	return &Cache{c}
}

type cache struct {
	sync.RWMutex
	items map[string]*Item
	option *Option
	optimizer *Optimizer
	size int
}

// Option option
type Option struct {
	ThresholdSize int // default is 0(unlimited)
	ThresholdAccess time.Duration // default is 0(not care)
	ThresholdAccessCount int64 // default is 0(not care)
}

// Set set item to cache.
// param key - key of item
// param value - value of item
// param expireIn - expire time
// return arg1 - Error
func (c *cache) Set(key string, value interface{}, expireIn time.Duration) error {
	supported, kind := c.IsSupported(value)
	if !supported {
		return errors.New("type of value is not supported. type = " + kind)
	}
	c.Lock()
	now := time.Now()
	var time time.Time
	if expireIn > 0 {
		time = now.Add(expireIn)
	} else {
		time = now.Add(DefaultExpiration)
	}
	item := &Item{key, value, 0, &time, 0, nil}
	item.Priority = c.priority(item)
	c.set(key, item)
	c.Unlock()
	return nil
}

func (c *cache) IsSupported(obj interface{}) (bool, string) {
	kind := reflect.TypeOf(obj).Kind()
	if kind == reflect.Chan || kind == reflect.Func {
		Warn("%s is unsupported.", kind);
		return false, kind.String();
	}
	
	return true, kind.String();
}

func (c *cache) SizeOfItem(item *Item) int {
	
	size := c.SizeOf(item.PriorityThan)
	size += c.SizeOf(item.AccessCount)
	size += c.SizeOf(item.Expiration)
	size += c.SizeOf(item.Key)
	size += c.SizeOf(item.LastAccess)
	size += c.SizeOf(item.Object)
	size += c.SizeOf(item.Priority)
	return size
}

func (c *cache) SizeOf(obj interface{}) int {
	t := reflect.TypeOf(obj)
	var o interface{}
	if t.Kind() == reflect.Ptr {
		v := reflect.ValueOf(obj)
		o = v.Elem()
	} else {
		o = obj
	}
	t = reflect.TypeOf(o)
	
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		v := reflect.ValueOf(o)
		elmSize := int(t.Elem().Size())
		elmNum := v.Len()
		return elmSize * elmNum
	case reflect.String:
		v := reflect.ValueOf(o)
		return v.Len()
	default:
		if t.Kind() == reflect.Struct {
			Debug("the case of a struct, size may not be correct. size is top level size of struct.");
		}
		return int(t.Size());
	}
}

func (c *cache) set(key string, item *Item) {
	beforeItemSize := 0
	if c.items[key] != nil {
		beforeItem := c.items[key]
		beforeItemSize = c.SizeOfItem(beforeItem)
	}
	c.items[key] = item
	c.size = c.size - beforeItemSize + c.SizeOfItem(item)
}

func (c *cache) Get(key string) (value *interface{}, found bool) {
	c.RLock()
	now := time.Now()
	item, found := c.get(key)
	if found {
		item.AccessCount = atomic.AddInt64(&item.AccessCount, 1)
		item.LastAccess = &now
		value = &item.Object
	}
	c.RUnlock()
	return value, found
}

// test case only
func (c *cache) GetItem(key string) (item *Item, found bool) {
	c.RLock()
	item, found = c.get(key)
	c.RUnlock()
	return item, found
}

// private function
func (c *cache) get(key string) (item *Item, found bool) {
	item = c.items[key]
	if item != nil {
		found = true
	}
	
	return item, found
}

func (c *cache) Del(key string) {
	c.Lock()
	c.del(key)
	c.Unlock()
}

func (c *cache) del(key string) {
	item, found := c.get(key)
	size := 0
	if found {
		size = c.SizeOfItem(item)
	}
	delete(c.items, key)
	//cache size
	c.size = c.size - size
}

// List of fineName in Cache
func (c *cache) List() []string {
	names := make([]string, len(c.items))

	i := 0
	for key := range c.items {
		names[i] = key
		i ++
	}

	return names
}

func (c *cache) Priority(key string) int {
	item, _ := c.GetItem(key)
	return c.priority(item)
}

func (c *cache) priority(item *Item) int {
	now := time.Now()
	priority := 0 //Items to be deleted
	if item == nil {
		return priority
	}
	
	// no expiration
	if item.Expiration == nil || item.Expiration.IsZero() {
		priority ++
	}
	// expiration > now
	if item.Expiration != nil && item.Expiration.After(now) {
		priority ++
	}
	// last access + threshold > now
	if c.option.ThresholdAccess != 0 && item.LastAccess != nil && item.LastAccess.Before(now) && item.LastAccess.Add(c.option.ThresholdAccess).After(now) {
		priority ++
	}
	// access count >= threshold
	if c.option.ThresholdAccessCount != 0 && item.AccessCount >= c.option.ThresholdAccessCount {
		priority ++
	}
	
	return priority
}

func (c *cache) Size() int {
	return c.size
}

// The optimized by priority
func (c *cache) Optimize() {
	Debug("before optimizing. files = %d size = %d bytes", len(c.items), c.size)

	// apply priority
	tmp := make(sortableItems, 0, len(c.items))
	for key, item := range c.items {
		priority := c.priority(item)
		if priority == 0 {
			c.Lock()
			c.del(key)
			Debug("optimizing delete key = %s", key)
			c.Unlock()
		} else {
			item.Priority = priority
			tmp = append(tmp, item)
		}
	}
	
	// sort
	sort.Sort(tmp)
	
	// compaction
	if 0 < c.option.ThresholdSize && c.size > c.option.ThresholdSize {
		for i, item := range tmp {
			if item != nil {
				c.Lock()
				key := item.Key
				tmp[i] = nil
				c.del(key)
				Debug("compaction delete key = %s", key)
				c.Unlock()
				if c.size <= c.option.ThresholdSize {
					break
				}
			} else {
				Warn("item is nil")
			}
		}
	}
	Debug("after optimizing. files = %d size = %d bytes", len(c.items), c.size)
}

// RunOptimizer run optimizing
// optimizer's task is ...
//  - update priority of item
//  - delete items of expired
// param interval - interval of optimize
func (c *cache) RunOptimizer(interval time.Duration) {
	if c.optimizer == nil {
		c.optimizer = &Optimizer{Interval: interval * time.Second}
	}
	
	go c.optimizer.Run(c)
}

// StopOptimizer stop optimizing
func (c *cache) StopOptimizer() {
	if c.optimizer != nil {
		c.optimizer.stop <- true
	}
}

// GetItems get item map from cache.
// return arg1 - item map
func (c *Cache) GetItems() map[string]*Item {
	return c.items
}

// Optimizer optimizer optimize cache.
type Optimizer struct {
	Interval time.Duration
	runing bool
	stop     chan bool
}

// Run run optimizer
// param c - instance of cache
func (o *Optimizer) Run(c *cache) {
	o.stop = make(chan bool)
	tick := time.Tick(o.Interval)
	for {
		select {
		case <-tick:
			if !o.runing {
				o.runing = true
				c.Optimize()
				o.runing = false
			}
			break
		case <-o.stop:
			return
		}
	}
}

// Item item of cache
type Item struct {
	Key string
	Object interface{}
	Priority int
	Expiration *time.Time
	AccessCount int64
	LastAccess *time.Time
}

// PriorityThan compare the priority
// param item - Item
// result arg1 - If this is higher than item, return true.
func (i *Item) PriorityThan(item *Item) bool {
	if i.Priority == item.Priority {
		if i.Expiration == item.Expiration {
			return true
		}
		return i.Expiration.After(*item.Expiration)
	}
	return (i.Priority) > (item.Priority)
}

type sortableItems []*Item
func (sitems sortableItems) Len() int           { return len(sitems) }
func (sitems sortableItems) Less(i, j int) bool {
	ii := sitems[i]
	ij := sitems[j]
	return !ii.PriorityThan(ij)
}
func (sitems sortableItems) Swap(i, j int)      { sitems[i], sitems[j] = sitems[j], sitems[i] }
