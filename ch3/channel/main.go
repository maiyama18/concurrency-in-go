package main

import "fmt"

func main() {
	channelOwner := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := 0; i < 5; i++ {
				fmt.Printf("sended %v\n", i)
				ch <- i
			}
		}()

		return ch
	}

	ch := channelOwner()
	for val := range ch {
		fmt.Printf("received %v\n", val)
	}
	fmt.Println("finished!")
}
