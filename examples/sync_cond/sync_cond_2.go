package sync_cond

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func requestStub(v int) int {
	time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	return v
}

/*
RunSyncCond2 Task is: need to do 100 "parallel" requests and
print "stub responses" from 1 to 100 in order without using data store for them
*/
func RunSyncCond2() {
	wg := sync.WaitGroup{}

	cond := sync.NewCond(&sync.Mutex{})
	numberResponseToShow := 1

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()

			response := requestStub(i)
			cond.L.Lock()
			for response != numberResponseToShow {
				cond.Wait()
			}
			fmt.Println(response)
			numberResponseToShow++
			cond.L.Unlock()
			cond.Broadcast()
		}()
	}

	wg.Wait()
}
