# go-cache
Library for caching

## Usage
```go 
package main

import "github.com/tico8/go-cache"

func main() {
  opt := cache.Option{
    ThresholdAccess: 0,
    ThresholdAccessCount: 0,
  }
  c := cache.New(opt)
  
  key := "testKey"
  value := "testValue"

  c.Set(key, value, time.Duration(10)) // expiration is 10 sec (default is 1 hour)
  resultValue, found := c.Get(key)
  c.Del(key)
}
```

## Option
* ThresholdSize  
Total size of cache
* ThresholdAccess  
If it is accessed within N seconds, priority +1
* ThresholdAccessCount  
If it is accessed more than N times, priority +1

## Optimizing of cache
If the size of cache is greater than the ThresholdSize value, it is possible to optimize caching .
Lower priority cache is deleted by optimizer.
```go 
package main

import "github.com/tico8/go-cache"

func main() {
  opt := cache.Option{
    ThresholdAccess: 1000000000, // 1 sec
    ThresholdAccessCount: 100, // 100 times
  }
  c := cache.New(opt)

  // Manually
  c.Optimize()

  // Executed in the 10 second intervals
  c.RunOptimizer(10)
	
  // Stop
  c.StopOptimizer()
}
```
