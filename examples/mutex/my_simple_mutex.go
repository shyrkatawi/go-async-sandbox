package mutex

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type mySimpleMutex struct {
	locked int32
}

func (m *mySimpleMutex) Lock() {
	for {
		if atomic.CompareAndSwapInt32(&m.locked, 0, 1) {
			return
		}
	}
}

func (m *mySimpleMutex) Unlock() {
	atomic.StoreInt32(&m.locked, 0)
}

func RunMyMutex() {
	wg := sync.WaitGroup{}
	m := mySimpleMutex{}

	counter := 0
	goroutinesNumber := 1000
	wg.Add(goroutinesNumber)

	for i := 0; i < goroutinesNumber; i++ {
		go func() {
			defer wg.Done()

			m.Lock()
			counter++
			m.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println(counter)
}
