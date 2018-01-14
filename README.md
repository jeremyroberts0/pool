# Pool
Pool is a thread pool library for Go.  It is designed to execute any number of jobs over a fix number of goroutines.

**It is highly recommend that you import from a git tag and not from master, the API may change**

## Usage
```go
package myApp

import (
  "github.com/jeremyroberts0/pool"
  "time"
)

func myJob() error {
  time.Sleep(time.Second)
  return nil
}

func doStuff() {
  pool := NewPool(5)

  pool.Add(myJob)
  pool.Add(myJob)
  pool.Add(myJob)
  pool.Add(myJob)
  pool.Add(myJob)

  errs := pool.Run()

  if len(errs) > 0 {
    fmt.Println("Errors in pool!")
  }
}
```