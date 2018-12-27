package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func genRand() interface{} {
	return rand.Intn(500000000)
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	repeatCh := make(chan interface{})
	go func() {
		defer close(repeatCh)
		for {
			select {
			case <-done:
				return
			case repeatCh <- fn():
			}
		}
	}()

	return repeatCh
}

func toInt(done <-chan interface{}, in <-chan interface{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- v.(int):
			}
		}
	}()

	return out
}

func take(done <-chan interface{}, in <-chan int, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case out <- <-in:
			}
		}
	}()

	return out
}

func primeFinder(done <-chan interface{}, in <-chan int) <-chan int {
	isPrime := func(n int) bool {
		for i := 2; i < n; i++ {
			if n%i == 0 {
				return false
			}
		}

		return true
	}

	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v := <-in:
				if isPrime(v) {
					out <- v
				}
			}
		}
	}()

	return out
}

func fanIn(done <-chan interface{}, ins []<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	for _, in := range ins {
		wg.Add(1)
		go func(in <-chan int) {
			defer wg.Done()
			for v := range in {
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntCh := toInt(done, repeatFn(done, genRand))

	primeChs := make([]<-chan int, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		primeChs[i] = primeFinder(done, randIntCh)
	}
	primeCh := take(done, fanIn(done, primeChs), 10)

	for n := range primeCh {
		fmt.Println(n)
	}

	fmt.Printf("completed. duration: %v\n", time.Since(start))
}
