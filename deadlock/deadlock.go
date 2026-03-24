package deadlock

import (
	"fmt"
	"sync"
	"time"
)

type valueWithMutex struct {
	name string
	v    int
	mu   sync.Mutex
}

func RunDeadLock() {
	var wg sync.WaitGroup

	var printSum = func(v1, v2 *valueWithMutex, index int) {
		defer wg.Done()

		fmt.Println(index, "trying to lock first", v1.name)
		v1.mu.Lock()
		fmt.Println(index, "locking first", v1.name)
		defer func() {
			fmt.Println(index, "unlocking first", v1.name)
			v1.mu.Unlock()
		}()

		time.Sleep(1 * time.Second)

		fmt.Println(index, "trying to lock second", v2.name)
		v2.mu.Lock()
		fmt.Println(index, "locking second", v2.name)
		defer func() {
			fmt.Println(index, "unlocking second", v2.name)
			v2.mu.Unlock()
		}()

		fmt.Println(index, v1.v+v2.v)
	}

	a := &valueWithMutex{name: "a"}
	b := &valueWithMutex{name: "b"}

	wg.Add(2)

	go printSum(a, b, 1)
	go printSum(b, a, 2)

	wg.Wait()
}
