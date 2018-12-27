package main

import (
	"fmt"
	"net/http"
)

type GetResult struct {
	Error error
	Response *http.Response
}

func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan GetResult {
		results := make(chan GetResult)
		go func() {
			defer close(results)
			for _, url := range urls {
				resp, err := http.Get(url)
				result := GetResult{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://non.existent.host", "https://yahoo.co.jp"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			continue
		}
		fmt.Printf("response: %v\n", result.Response.Status)
	}
}
