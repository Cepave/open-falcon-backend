package simplecache

import (
	"bytes"
	"os"
	"runtime"
	"sync"
	"testing"
)

type TestStruct struct {
	Num  int
	Name string
	Data interface{}
}

func newTestStruct(num int, name string) *TestStruct {
	return &TestStruct{Num: num, Name: name, Data: nil}
}

func (t *TestStruct) equal(other *TestStruct) bool {
	return (t.Name == other.Name && t.Num == other.Num && t.Data == other.Data)
}

func contain(array []string, key string) bool {
	for _, item := range array {
		if item == key {
			return true
		}
	}
	return false
}

// Unit Test
func TestCache(t *testing.T) {
	cache := NewCache()
	a := newTestStruct(1, "a")
	b := newTestStruct(2, "b")
	a1 := newTestStruct(-1, "a1")
	a1.Data = a

	// Exist
	if cache.Exist("a") {
		t.Error("Exist")
	}
	// Get
	item, found := cache.Get("a")
	if !(item == nil && !found) {
		t.Error("Get")
	}
	// Set
	// SetIfNonExistent
	cache.Set("a", a)
	item, found = cache.Get("a")
	if !(found && item != nil) {
		t.Error("Set")
	}
	titem := item.(*TestStruct)
	if titem.Name != "a" || titem.Num != 1 {
		t.Error("Set")
	}

	cache.SetIfNonExistent("a", a1)
	item, found = cache.Get("a")
	titem = item.(*TestStruct)
	if !(titem.Name == "a" && titem.Num == 1) {
		t.Error("SetIfNonExistent")
	}
	cache.SetIfNonExistent("b", b)
	item, found = cache.Get("b")
	titem = item.(*TestStruct)
	if !(titem.Name == "b" && titem.Num == 2) {
		t.Error("SetIfNonExistent")
	}
	// Keys
	keys := cache.Keys()
	if !(len(keys) == 2 && contain(keys, "a") && contain(keys, "b")) {
		t.Error("Keys")
	}
	// Len
	length := cache.Len()
	if !(length == 2) {
		t.Error("Len")
	}
	// Remove
	cache.Remove("b")
	if !(cache.Exist("a") && !cache.Exist("b")) {
		t.Error("Remove")
	}
	cache.Remove("a")
	if !(!cache.Exist("a") && !cache.Exist("b")) {
		t.Error("Remove")
	}
	// RemoveAll
	cache.Set("a", a)
	cache.Set("b", b)
	cache.RemoveAll()
	if !(cache.Len() == 0) {
		t.Error("RemoveAll")
	}
	// Save && Load
	cache.RemoveAll()
	cache.Set("a", a)
	cache.Set("b", b)
	fp := &bytes.Buffer{}
	if !(cache.Save(fp) == nil) {
		t.Error("Save")
	}
	cache.RemoveAll()
	cache.Set("a", a1)
	if !(cache.Load(fp) == nil) {
		t.Error("Load")
	}
	if !(cache.Len() == 2 && cache.Exist("a") && cache.Exist("b")) {
		t.Error("Save && Load")
	}
	item, _ = cache.Get("a")
	titem = item.(*TestStruct)
	if !(titem.Name == "a1" && titem.Num == -1 && titem.Data.(*TestStruct).equal(a)) {
		t.Error("Save && Load")
	}
	// SaveToFile && LoadFromFile
	cache.RemoveAll()
	cache.Set("a", a)
	cache.Set("b", b)
	cache.SaveToFile("./cache.file")
	cache.RemoveAll()
	cache.Set("a", a1)
	cache.LoadFromFile("./cache.file")
	os.Remove("./cache.file")
	if !(cache.Len() == 2 && cache.Exist("a") && cache.Exist("b")) {
		t.Error("SaveToFile && LoadFromFile")
	}
	item, _ = cache.Get("a")
	titem = item.(*TestStruct)
	if !(titem.Name == "a1" && titem.Num == -1 && titem.Data.(*TestStruct).equal(a)) {
		t.Error("SaveToFile && LoadFromFile")
	}
	// EchoVsn
	if !(VSN == EchoVsn()) {
		t.Error("EchoVsn")
	}
}

// Benchmark Test
func BenchmarkCacheGet(b *testing.B) {
	b.StopTimer()
	cache := NewCache()
	cache.Set("foo", "bar")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("foo")
	}
}

func BenchmarkCacheGetConcurrent(b *testing.B) {
	b.StopTimer()
	cache := NewCache()
	cache.Set("foo", "bar")
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for j := 0; j < workers; j++ {
		go func() {
			for i := 0; i < each; i++ {
				cache.Get("foo")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkMapGet(b *testing.B) {
	b.StopTimer()
	cache := map[string]string{"foo": "bar"}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_, _ = cache["foo"]
		mu.RUnlock()
	}
}

func BenchmarkMapGetConcurrent(b *testing.B) {
	b.StopTimer()
	cache := map[string]string{"foo": "bar"}
	mu := sync.RWMutex{}
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for j := 0; j < workers; j++ {
		go func() {
			for i := 0; i < each; i++ {
				mu.RLock()
				_, _ = cache["foo"]
				mu.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCacheSet(b *testing.B) {
	b.StopTimer()
	cache := NewCache()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("foo", "bar")
	}
}

func BenchmarkCacheSetConcurrent(b *testing.B) {
	b.StopTimer()
	cache := NewCache()
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for j := 0; j < workers; j++ {
		go func() {
			for i := 0; i < each; i++ {
				cache.Set("foo", "bar")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkMapSet(b *testing.B) {
	b.StopTimer()
	cache := make(map[string]string)
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		cache["foo"] = "bar"
		mu.Unlock()
	}
}

func BenchmarkMapSetConcurrent(b *testing.B) {
	b.StopTimer()
	cache := make(map[string]string)
	mu := sync.RWMutex{}
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for j := 0; j < workers; j++ {
		go func() {
			for i := 0; i < each; i++ {
				mu.Lock()
				cache["foo"] = "bar"
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCacheSetIfNonExistent(b *testing.B) {
	b.StopTimer()
	cache := NewCache()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.SetIfNonExistent("foo", "bar")
	}
}

func BenchmarkCacheSetRemove(b *testing.B) {
	b.StopTimer()
	cache := NewCache()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("foo", "bar")
		cache.Remove("foo")
	}
}
func BenchmarkMapSetRemove(b *testing.B) {
	b.StopTimer()
	cache := make(map[string]string)
	mu := new(sync.RWMutex)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		cache["foo"] = "bar"
		mu.Unlock()
		mu.Lock()
		delete(cache, "foo")
		mu.Unlock()
	}
}
