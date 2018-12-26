package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)

			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()

		return terminated
	}

	done := make(chan interface{})
	fmt.Println("launching doWork...")
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("canceling doWork...")
		close(done)
	}()

	<-terminated
	fmt.Println("done")
}
