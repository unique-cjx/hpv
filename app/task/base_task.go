package task

import (
	jsoniter "github.com/json-iterator/go"
	"sync"
)

var (
	DepartChan    chan *DepartRows
	DepartStorage *departStorage
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
)

type departStorage struct {
	Lock sync.RWMutex
	Dids []int64
}

func init() {
	DepartChan = make(chan *DepartRows, 2<<4)
	DepartStorage = &departStorage{Dids: make([]int64, 0)}
}
