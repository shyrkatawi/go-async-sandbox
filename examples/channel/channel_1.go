package channel

import (
	"fmt"
	"sync"
)

func RunChannel1() {
	/*
		counter := 0
		for i := 0; i < 1000; i++ {
			go func() {
				counter++
			}()
		}
		fmt.Println(counter) // must show 1000 using channel to solution
	*/

	wg := sync.WaitGroup{}
	ch := make(chan any, 1)

	counter := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ch <- nil
			counter++
			<-ch
		}()
	}

	wg.Wait()
	fmt.Println(counter)
}
