package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
	"strconv"
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
	return &LockKey{
		locks:make(map[string]*lock),
		mutex:&sync.Mutex{},
	}
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

func main() {
	var wg sync.WaitGroup
	keys := NewKeyLock()
	for i:=0; i< 100; i++{
		num := i
		wg.Add(1)
		go func(wg *sync.WaitGroup, id int){
			for i:= 0; i < 10; i++{
				keys.Block(strconv.Itoa(rand.Intn(10)))
				time.Sleep(time.Second * time.Duration(rand.Intn(3)))
				keys.Unblock(strconv.Itoa(rand.Intn(10)))
				fmt.Printf("%d is done", id)
				time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			}
			wg.Done()
		}(&wg, num)
	}
	wg.Wait()
}