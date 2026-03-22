package main

import "go-async-sandbox/examples/sync_cond"

func main() {
	//deadlock.RunDeadLock()
	//livelock.RunLivelock()
	//sync_cond.RunSyncCond1()
	//starvation.RunStarvation()
	//mutex.RunMyMutex()
	sync_cond.RunSyncCond2()
}
