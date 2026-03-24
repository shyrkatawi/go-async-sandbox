package livelock

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func RunLivelock() {
	cadence := sync.NewCond(&sync.Mutex{})

	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDirection := func(directionName string, tries *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, " %v", directionName)
		atomic.AddInt32(tries, 1)
		takeStep()

		if atomic.LoadInt32(tries) == 1 {
			fmt.Println(out, "Success")
			return true
		}

		takeStep()
		atomic.AddInt32(tries, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDirection("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDirection("right", &right, out) }

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer

		defer func() {
			fmt.Println(out.String())
		}()

		defer walking.Done()

		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	peopleInHallwayWg := sync.WaitGroup{}
	peopleInHallwayWg.Add(2)

	go walk(&peopleInHallwayWg, "Alice")
	go walk(&peopleInHallwayWg, "Barbara")

	peopleInHallwayWg.Wait()
}
