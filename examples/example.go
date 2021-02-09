package main

import (
	"fmt"

	"github.com/vitorsalgado/gow/pkg/worker"
)

func main() {
	c := make(chan bool)

	dispatcher := worker.NewDispatcher(10)
	dispatcher.Run()

	dispatcher.Dispatch(func() (id string, err error) {
		id = "job-id-#1"
		fmt.Printf("Running job: %v\n", id)
		return id, nil
	})

	dispatcher.Dispatch(func() (id string, err error) {
		id = "job-id-#2"
		fmt.Printf("Running job: %v\n", id)
		return id, nil
	})

	dispatcher.Quit()

	<-c
}
