package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup // for waiting
	var mu sync.Mutex
	counter := 0 // the shared number

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++ // all 1000 goroutines touch the same variable
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println(counter) // expected 1000, but often less
}
