package channel

import "fmt"

func RunChannel3() {
	// explain the code step by step
	type c chan c

	ch := make(chan c, 1)
	ch <- ch

	for i := 0; i < 100; i++ {
		select {
		case <-ch:
		case <-ch:
			ch <- ch
		default:
			fmt.Println(i)
			return
		}
	}
}
