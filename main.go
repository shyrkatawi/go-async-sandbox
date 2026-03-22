package main

import (
	"go-async-sandbox/examples/mutex"
)

func main() {
	//deadlock.RunDeadLock()
	//livelock.RunLivelock()
	//sync_cond.RunSyncCond()
	//starvation.RunStarvation()
	mutex.RunMyMutex()
}
