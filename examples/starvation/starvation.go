package starvation

import (
	"fmt"
	"sync"
	"time"
)

func RunStarvation() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	const timeToRun = 1 * time.Second

	greedyWorker := func() {
		defer wg.Done()

		numberOfOperations := 0
		for begin := time.Now(); time.Since(begin) <= timeToRun; {
			mu.Lock()
			time.Sleep(3 * time.Nanosecond)
			mu.Unlock()

			numberOfOperations++
		}

		fmt.Println("greedyWorker did", numberOfOperations, "operations")
	}

	politeWorker := func() {
		defer wg.Done()

		numberOfOperations := 0
		for begin := time.Now(); time.Since(begin) <= timeToRun; {
			for i := 0; i < 3; i++ {
				mu.Lock()
				time.Sleep(time.Nanosecond)
				mu.Unlock()
			}

			numberOfOperations++
		}

		fmt.Println("greedyWorker did", numberOfOperations, "operations")
	}

	wg.Add(2)
	greedyWorker()
	politeWorker()
	wg.Wait()
}
