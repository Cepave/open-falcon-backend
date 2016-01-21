package simplecache

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	VSN = "0.0.1"
)

// basic cache data item
type Item struct {
	Object interface{}
}

func NewItem(obj interface{}) *Item {
	return &Item{Object: obj}
}

// cache object
type Cache struct {
	sync.RWMutex
	items map[string]*Item
}

func NewCache() *Cache {
	return &Cache{items: make(map[string]*Item)}
}

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	obj, found := c.get(key)
	c.RUnlock()
	return obj, found
}
func (c *Cache) get(key string) (interface{}, bool) {
	item, found := c.items[key]
	if found {
		return item.Object, true
	} else {
		return nil, false
	}
}

// Add an item to the cache, replacing any existing item.
func (c *Cache) Set(key string, val interface{}) {
	c.Lock()
	c.set(key, val)
	c.Unlock()
}

// Add an item to the cache if no item with the same key exists
func (c *Cache) SetIfNonExistent(key string, val interface{}) {
	c.Lock()
	if !c.exist(key) {
		c.set(key, val)
	}
	c.Unlock()
}
func (c *Cache) set(key string, val interface{}) {
	c.items[key] = NewItem(val)
}

// Return whether or not item with the $key exits in the cache
func (c *Cache) Exist(key string) bool {
	c.RLock()
	found := c.exist(key)
	c.RUnlock()
	return found
}
func (c *Cache) exist(key string) bool {
	_, found := c.items[key]
	return found
}

// Remove an item from the cache. Does nothing if the key is not in the cache.
func (c *Cache) Remove(key string) {
	c.Lock()
	delete(c.items, key)
	c.Unlock()
}

// Remove all item from the cache
func (c *Cache) RemoveAll() {
	c.Lock()
	defer c.Unlock()

	c.items = make(map[string]*Item)
}

// Return keys of all the items.
func (c *Cache) Keys() []string {
	c.RLock()
	defer c.RUnlock()
	return c.keys()
}
func (c *Cache) keys() []string {
	keys := make([]string, 0, c.len())
	for key, _ := range c.items {
		keys = append(keys, key)
	}
	return keys
}

// Returns the number of items in the cache.
func (c *Cache) Len() int {
	c.RLock()
	defer c.RUnlock()

	return c.len()
}
func (c *Cache) len() int {
	return len(c.items)
}

// Save the cache's items to the given filename, creating the file if it
// doesn't exist, and overwriting it if it does.
func (c *Cache) SaveToFile(filename string) error {
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = c.Save(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}
func (c *Cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()

	c.RLock()
	keys := c.keys()
	c.RUnlock()
	// Break one long rlock into c.Len short ones
	for _, k := range keys {
		v, found := c.Get(k)
		if found { // not thread safe here, so we have to check it
			gob.Register(v)
		}
	}
	err = enc.Encode(&c.items)
	return
}

// Load and add cache items from the given filename, excluding any items with
// keys that already exist in the current cache.
func (c *Cache) LoadFromFile(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	err = c.Load(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}
func (c *Cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := make(map[string]*Item)
	err := dec.Decode(&items)
	if err != nil {
		return err
	}

	for k, v := range items {
		// use len(items) short locks instead of one long lock
		c.SetIfNonExistent(k, v)
	}
	return nil
}

// Show Version of cache source code.
func EchoVsn() string {
	fmt.Println(VSN)
	return VSN
}
