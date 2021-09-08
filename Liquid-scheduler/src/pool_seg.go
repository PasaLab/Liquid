package main

import "sync"

type PoolSeg struct {
	ID        int
	Nodes     map[string]*NodeStatus
	Lock      sync.Mutex
	Next      *PoolSeg
}
