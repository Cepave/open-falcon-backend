# cache
an in-memory key:value cache library for Go, suitable for single-machine applications.<br />
<b>features:</b><br />
+ multi-threads safe in most cases. the cache can be safely used by multiple goroutines
+ persistent storage. the entire cache can be saved to and loaded from a file to recover from downtime quickly

## installation
    git clone https://github.com/niean/cache
    go get https://github.com/niean/cache

## peference
```
```

## usage
```go
	import (
		"fmt"
		cache "github.com/niean/gotools/cache/simplecache"
		"time"
	)

	func main() {
		// simplecache
		c := cache.NewCache()
		c.Set("foo", "bar")

		foo, found := c.Get("foo")
		if found {
			fmt.Println(foo)
		}

		v, found := c.Get("foo")
		if found {
			fooString := v.(string)
			fmt.Printf("%s", fooString)
		}

		if v, found := c.Get("foo"); found {
			fooString := v.(string)
			fmt.Printf("%s", fooString)
		}

		var fooString string
		if v, found := c.Get("foo"); found {
			fooString = v.(string)
		}

		foo := &MyStruct{Num: 1}
		c.Set("foo", foo)
		x, _ := c.Get("foo")
		foo := x.(*MyStruct)
		fmt.Println(foo.Num)
		foo.Num++
		x, _ := c.Get("foo")
		foo := x.(*MyStruct)
		foo.Println(foo.Num)
		// will print:
		// 1
		// 2

		// timedcache
		// TODO
	}
```
## reference
simplified from [go-cache](https://github.com/pmylund/go-cache)
