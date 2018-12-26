package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	newRandCh := func(done <-chan interface{}) <-chan int {
		randCh := make(chan int)
		go func() {
			defer fmt.Println("newRandCh exited")
			defer close(randCh)

			for {
				select {
					case randCh <- rand.Int():
					case <-done:
						return
				}
			}
		}()

		return randCh
	}

	done := make(chan interface{})
	randCh := newRandCh(done)
	for i := 0; i < 3; i++ {
		fmt.Printf("rand %v: %d\n", i, <-randCh)
	}
	close(done)
	time.Sleep(500 * time.Millisecond)
}
