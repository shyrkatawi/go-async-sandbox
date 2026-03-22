package main

import "go-async-sandbox/examples/starvation"

func main() {
	//deadlock.RunDeadLock()
	//livelock.RunLivelock()
	//sync_cond.RunSyncCond()
	starvation.RunStarvation()
}
