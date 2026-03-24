package channel

import (
	"fmt"
	"sync"
	"time"
)

/*
	implement the function
	func f(channels ...<-chan int) <-chan int {}
*/

func withForLoop() {
	f := func(channels ...<-chan int) <-chan int {
		returnCn := make(chan int)

		closedChNumber := 0

		go func() {
			for i := 0; i < len(channels); i++ {
				go func() {
					for {
						v, ok := <-channels[i]
						if !ok {
							closedChNumber++
							if closedChNumber == len(channels) {
								close(returnCn)
							}
							return
						}
						returnCn <- v
					}
				}()
			}
		}()

		return returnCn
	}

	chArr := make([]<-chan int, 10)
	for i := 0; i < len(chArr); i++ {
		ch := make(chan int)
		chArr[i] = ch

		go func() {
			ch <- i
			time.Sleep(time.Second)

			ch <- i * 100
			time.Sleep(time.Second)

			close(ch)
		}()
	}

	ch := f(chArr...)
	for {
		v, ok := <-ch
		fmt.Println("reading from result channel:", v)
		if !ok {
			fmt.Println("result channel is closed")
			break
		}
	}
}

func withRange() {
	f := func(channels ...<-chan int) <-chan int {
		resCh := make(chan int)

		wg := sync.WaitGroup{}

		go func() {
			for _, ch := range channels {
				wg.Add(1)
				go func() {
					defer wg.Done()

					for v := range ch {
						resCh <- v
					}
				}()
			}

			wg.Wait()
			close(resCh)
		}()

		return resCh
	}

	chArr := make([]<-chan int, 10)

	for i := range chArr {
		ch := make(chan int)
		chArr[i] = ch
		go func() {
			ch <- i
			time.Sleep(time.Second)

			ch <- i * 100
			time.Sleep(time.Second)

			close(ch)
		}()
	}

	resCh := f(chArr...)

	for v := range resCh {
		fmt.Println("reading from result channel:", v)
	}
	fmt.Println("result channel is closed")
}

func RunChannel4() {
	withForLoop()
	withRange()
}
