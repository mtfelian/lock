package lock

import "sync"

// Locker is an interface for locking
type Locker interface {
	Block()
	Unblock()
}

// KeyLocker is an interface to lock by string key
type KeyLocker interface {
	Locker
	NewKey()
}

// lock is a mutex lock by key object
type lock struct {
	mutex *sync.Mutex
	once  *sync.Once
}

// NewLock creates new lock object
func NewLock() *lock {
	return &lock{&sync.Mutex{}, &sync.Once{}}
}

// Block hangs is obj already blocked
func (obj *lock) Block() {
	obj.mutex.Lock()
	obj.once = &sync.Once{}
}

// Unblock unblocks obj
func (obj *lock) Unblock() {
	obj.once.Do(func() { obj.mutex.Unlock() })
}

// keyLock is a mutex lock by key object
type keyLock struct {
	sync.RWMutex // this is a sync for map of mutexes below
	mutex        map[string]*sync.Mutex
}

// NewKeyLock creates new keyLock object
func NewKeyLock() *keyLock {
	return &keyLock{
		mutex: make(map[string]*sync.Mutex),
	}
}

// newKey creates new keyLock object by key
func (lk *keyLock) newKey(key string) {
	lk.mutex[key] = &sync.Mutex{}
}

// Block blocks lk. If called more than once waits for an unlock
func (lk *keyLock) Block(key string) {
	lk.RLock()
	m, ok := lk.mutex[key]
	lk.RUnlock()

	if !ok {
		lk.Lock()
		lk.newKey(key)
		m, _ = lk.mutex[key]
		lk.Unlock()
	}

	m.Lock()
}

// Unblock unblocks lk by key
func (lk *keyLock) Unblock(key string) {
	lk.RLock()
	m, ok := lk.mutex[key]
	lk.RUnlock()
	if ok {
		m.Unlock()
	}
}
