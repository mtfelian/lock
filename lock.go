package lock

import "sync"

// KeyLock Mutex lock by key interface
type KeyLock struct {
	sync.RWMutex // этот мьютекс защищает набор мьютексов ниже
	mutex        map[string]*sync.Mutex
}

// Block Блокировка с ожиданием по ключу
func (kl *KeyLock) Block(key string) {

	kl.RLock()
	m, ok := kl.mutex[key]
	kl.RUnlock()
	if !ok {
		kl.Lock()
		kl.mutex[key] = &sync.Mutex{}
		m, ok = kl.mutex[key]
		kl.Unlock()
		if !ok {
			panic("lock.block(): can't create mutex")
		}
	}

	m.Lock()
}

// Unblock Раблокировка по ключу
func (kl *KeyLock) Unblock(key string) {

	kl.RLock()
	m, ok := kl.mutex[key]
	kl.RUnlock()
	if ok {
		m.Unlock()
		// это можешь убрать, если не хочешь удалять созданные мьютексы
		kl.Lock()
		delete(kl.mutex, key)
		kl.Unlock()
	}
}

// NewKeyLock конструктор
func NewKeyLock() *KeyLock {

	return &KeyLock{
		mutex: make(map[string]*sync.Mutex),
	}
}
