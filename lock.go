package lock

import "sync"

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

// LockKey Mutex lock by key
type LockKey struct {
	sync.RWMutex // этот мьютекс защищает набор мьютексов ниже
	mutex        map[string]*sync.Mutex
}

// NewKeyLock конструктор
func NewKeyLock() *LockKey {

	return &LockKey{
		mutex: make(map[string]*sync.Mutex),
	}
}

func (lk *LockKey) newKey(key string) {

	lk.mutex[key] = &sync.Mutex{}
}

// Block вызывает блокировку.
// При повторном вызове подвисает в ожидании разблокировки
func (lk *LockKey) Block(key string) {

	lk.RLock()
	m, ok := lk.mutex[key]
	lk.RUnlock()
	if !ok {
		lk.Lock()
		lk.mutex[key] = &sync.Mutex{}
		m, ok = lk.mutex[key]
		lk.Unlock()
		if !ok {
			panic("lock.block(): can't create mutex")
		}
	}

	m.Lock()
}

// Unblock вызывает разблокировку.
// При повторном вызове ничего не происходит.
func (lk *LockKey) Unblock(key string) {

	lk.RLock()
	m, ok := lk.mutex[key]
	lk.RUnlock()
	if ok {
		m.Unlock()
	}
}
