package constants

import "sync"

var (
	Mutex    = sync.Mutex{}
	MutexMap = map[int]*sync.Mutex{}
)
