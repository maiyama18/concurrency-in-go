package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func main() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1)
		}
	}

	consumer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, l sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)

		startTime := time.Now()
		go producer(&wg, l)
		for i := count; i > 0; i-- {
			go consumer(&wg, l)
		}

		wg.Wait()
		return time.Since(startTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var rwmtx sync.RWMutex
	var mtx sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))

		timeForMutex := test(count, &mtx)
		timeForRWMutex := test(count, &rwmtx)
		fmt.Fprintf(tw, "%v\t%v\t%v\n", count, timeForRWMutex, timeForMutex)
	}
}
