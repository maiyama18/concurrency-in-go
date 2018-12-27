package main

import "fmt"

func main() {
	generator := func(done <-chan interface{}, ints ...int) <-chan int {
		intCh := make(chan int)
		go func() {
			defer close(intCh)

			for _, i := range ints {
				select {
				case <-done:
					return
				case intCh <- i:
				}
			}
		}()

		return intCh
	}

	add := func(done <-chan interface{}, inCh <-chan int, additive int) <-chan int {
		outCh := make(chan int)
		go func() {
			defer close(outCh)

			for i := range inCh {
				select {
				case <-done:
					return
				case outCh <- i + additive:
				}
			}
		}()

		return outCh
	}

	mul := func(done <-chan interface{}, inCh <-chan int, multiplier int) <-chan int {
		outCh := make(chan int)
		go func() {
			defer close(outCh)

			for i := range inCh {
				select {
				case <-done:
					return
				case outCh <- i * multiplier:
				}
			}
		}()

		return outCh
	}

	done := make(chan interface{})
	defer close(done)

	intCh := generator(done, 1, 2, 3, 4)
	for v := range mul(done, add(done, intCh, 3), 2) {
		fmt.Println(v)
	}
}
