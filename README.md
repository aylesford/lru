# lru
This package of `lru` is the most concise and efficient LRU algorithm based on golang.

# Example
Quick start:
```go
package main

import (
	"fmt"

	"github.com/aylesford/lru"
)

func main() {
	cache := lru.NewLRU()

	cache.Add("aa", "AA")
	cache.Add("bb", "BB")
	cache.Add("cc", "CC")

	val, ok := cache.Get("cc")
	if !ok {
		fmt.Println("cache missing")
		return
	}

	str, ok := val.(string)
	if !ok {
		fmt.Println("type error")
		return
	}

	fmt.Println("result: ", str)
}
```
