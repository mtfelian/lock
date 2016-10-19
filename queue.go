package lock

import (
	"errors"
	"sync"
)

type SyncedQueue struct {
	queue []interface{}
	sync.Mutex
}

func NewSyncedQueue() SyncedQueue {
	return SyncedQueue{
		[]interface{}{},
		sync.Mutex{},
	}
}

// кладёт объект в очередь
func (q *SyncedQueue) Push(object interface{}) {
	q.queue = append(q.queue, object)
}

// возвращает размер очереди
func (q *SyncedQueue) Len() int {
	return len(q.queue)
}

// достаёт из очереди объект
func (q *SyncedQueue) Pop() (interface{}, error) {
	if q.Len() == 0 {
		return nil, errors.New("Очередь пуста")
	}
	poppedObject := q.queue[0]
	q.queue = q.queue[1:]
	return poppedObject, nil
}
