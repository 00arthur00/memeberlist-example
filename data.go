package main

import (
	"sync"

	"github.com/hashicorp/memberlist"
)

var (
	mtx          sync.Mutex
	data         map[string]string
	limitedQueue memberlist.TransmitLimitedQueue
)
