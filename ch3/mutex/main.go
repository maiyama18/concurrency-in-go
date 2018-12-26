package main

import (
	"fmt"
	"sync"
)

const times = 5

func main() {
	var count int
	var mtx sync.Mutex

	increment := func() {
		mtx.Lock()
		defer mtx.Unlock()
		count++
		fmt.Printf("incrementing: %v\n", count)
	}

	decrement := func() {
		mtx.Lock()
		defer mtx.Unlock()
		count--
		fmt.Printf("decrementing: %v\n", count)
	}

	var wg sync.WaitGroup

	for i := 0; i < times; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			increment()
		}()
	}
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			decrement()
		}()
	}

	wg.Wait()
}
