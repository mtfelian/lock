package lock

import (
	"sync"
)

// интерфейс блокировки
type Locker interface {
	Block()
	Unblock()
}

// интерфейс блокировки по строковому ключу
type KeyLocker interface {
	Locker
	NewKey()
}

type lock struct {
	mutex *sync.Mutex
	once  *sync.Once
}

func NewLock() *lock {
	return &lock{&sync.Mutex{}, &sync.Once{}}
}

func (obj *lock) Block() {
	obj.mutex.Lock()
	obj.once = &sync.Once{}
}

func (obj *lock) Unblock() {
	obj.once.Do(func() { obj.mutex.Unlock() })
}

type LockKey struct {
	mutex *sync.Mutex
	locks map[string]*lock
}

func NewKeyLock() *LockKey {
	return &LockKey{make(map[string]*lock}
}

func (obj *LockKey) newKey(key string) {
	obj.locks[key] = NewLock()
}

// Block вызывает блокировку.
// При повторном вызове подвисает в ожидании разблокировки
func (obj *LockKey) Block(key string) {
	obj.mutex.Lock()
	lock, ok := obj.locks[key]
	if !ok {
		obj.newKey(key)
		lock = obj.locks[key]
	}
	obj.mutex.Unlock()
	lock.Block()
}

// Unblock вызывает разблокировку.
// При повторном вызове ничего не происходит.
func (obj *LockKey) Unblock(key string) {
	lock, ok := obj.locks[key]
	if !ok {
		return
	}
	lock.Unblock()
}
