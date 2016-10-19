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
	mutex map[string]*sync.Mutex
	once  map[string]*sync.Once
}

func NewKeyLock() *LockKey {
	return &LockKey{make(map[string]*sync.Mutex), make(map[string]*sync.Once)}
}

func (obj *LockKey) newKey(key string) {
	obj.mutex[key] = &sync.Mutex{}
	obj.once[key] = &sync.Once{}
}

// Block вызывает блокировку.
// При повторном вызове подвисает в ожидании разблокировки
func (obj *LockKey) Block(key string) {
	m, ok := obj.mutex[key]
	if !ok {
		obj.newKey(key)
		m = obj.mutex[key]
	}
	m.Lock()
	obj.once[key] = &sync.Once{}
}

// Unblock вызывает разблокировку.
// При повторном вызове ничего не происходит.
func (obj *LockKey) Unblock(key string) {
	m, ok := obj.mutex[key]
	if !ok {
		return
	}
	obj.once[key].Do(func() { m.Unlock() })
}
