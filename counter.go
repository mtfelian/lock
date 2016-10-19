package lock

import (
	"sync"
)

type SyncedCounter struct {
	count int
	sync.Mutex
}

func NewSyncedCounter(initialValue int) SyncedCounter {
	return SyncedCounter{
		initialValue,
		sync.Mutex{},
	}
}

// увеличивает значение счётчика на 1
func (c *SyncedCounter) Inc() {
	c.count++
}

// уменьшает значение счётчика на 1
func (c *SyncedCounter) Dec() {
	c.count--
}

// возвращает значение счётчика
func (c *SyncedCounter) Get() int {
	return c.count
}
