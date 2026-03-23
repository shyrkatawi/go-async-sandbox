package channel

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func RunChannel2() {
	/*
			  run 10 slowFunction in "parallel" with each function timeout 3 sec
			  show result of execution in console

		    // Do not modify slowFunction
				func slowFunction() int {
					n:= rand.Intn(10)
					time.Sleep(time.Second * time.Duration(n))
					return n
				}

				func runner() {
					for i := 0; i < 10; i++ {
			    	result := slowFunction()
						fmt.Println(result)
					}
				}
	*/

	wg := sync.WaitGroup{}

	slowFunction := func() int {
		n := rand.Intn(10)
		time.Sleep(time.Second * time.Duration(n))
		return n
	}

	runner := func() {
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(3 * time.Second)
			cancel()
		}()

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resultCh := make(chan int)

				go func() {
					resultCh <- slowFunction()
					close(resultCh)
				}()

				select {
				case result := <-resultCh:
					fmt.Println(i, "result is", result)
				case <-ctx.Done():
					fmt.Println(i, "timeout")
				}
			}()
		}
	}

	runner()

	wg.Wait()
}
