/*
=== Example usage of Cache ===

	type Xyz struct {
		// ...
	}

	type XyzCache struct {
		cache.Store[Xyz]
	}

	func (c *XyzCache) Update(input any) error {
		c.Lock.Lock()
		defer c.Lock.Unlock()

		// Get Data
		c.Data = Xyz{
			// ...
		}

		cacheValidPeriod, _ := time.ParseDuration("12h")
		c.ExpiryTime = time.Now().Add(cacheValidPeriod)
		return nil
	}

var xyzCache cache.Cache[Xyz] = &XyzCache{}

func main() {

		input := map[string]any{
			"key": "value",
		}
		data, err := xyzCache.GetData(xyzCache, input)
		if err != nil {
			fmt.Println("Error while accessing xyzCache : ", err)
			return
		}
		fmt.Println(data)
	}
*/
package cache

import (
	"sync"
	"time"
)

type Store[T any] struct {
	Data       T
	ExpiryTime time.Time
	Lock       sync.RWMutex
}

type Cache[T any] interface {
	IsValid() bool
	GetData(Cache[T], any) (T, error)

	// Implement update per cache basis
	Update(any) error
}

func (c *Store[T]) IsValid() bool {
	if c.ExpiryTime.IsZero() || c.ExpiryTime.Before(time.Now()) {
		return false
	}
	return true
}

func (c *Store[T]) GetData(ci Cache[T], input any) (T, error) {

	if !c.IsValid() {
		err := ci.Update(input)
		if err != nil {
			return c.Data, err
		}
	}

	c.Lock.RLock()
	defer c.Lock.RUnlock()
	return c.Data, nil
}
