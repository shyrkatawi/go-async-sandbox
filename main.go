package main

import (
	"go-async-sandbox/examples/livelock"
)

func main() {
	//deadlock.RunDeadLock()
	livelock.RunLivelock()
	//sync_cond.RunSyncCond()
}
