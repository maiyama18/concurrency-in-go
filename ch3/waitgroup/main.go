package main

import (
	"fmt"
	"sync"
)

const numGreeters = 5

func main() {
	var wg sync.WaitGroup

	for i := 0; i < numGreeters; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fmt.Printf("hello from %v\n", index)
		}(i)
	}

	wg.Wait()
}


