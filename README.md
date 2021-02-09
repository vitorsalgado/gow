# Gow

![ci](https://github.com/vitorsalgado/gow/workflows/ci/badge.svg)

Simple Golang worker implementation.  
This was a playground to experiment **goroutines** and **channels**.

## Usage
```
package main

import "github.com/vitorsalgado/gow/pkg/worker"

func main() {
	dispatcher := worker.NewDispatcher(10)
	dispatcher.Run()
	dispatcher.Dispatch(func() (id string, err error) { })
}
```

