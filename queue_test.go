package lock

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func getQueue() SyncedQueue {
	q := NewSyncedQueue()
	q.queue = []interface{}{1, 2, 3}
	return q
}

func TestQueuePushToNonEmpty(t *testing.T) {
	q := getQueue()
	q.Push(4)
	if q.Len() != 4 {
		t.Fatalf("Ожидалась длина %v, а получено: %v", 4, q.Len())
	}
	lastElement := q.queue[3].(int)
	if lastElement != 4 {
		t.Fatalf("Ожидалось получить %v, получено: %v", 4, lastElement)
	}
}

func TestQueuePushToEmpty(t *testing.T) {
	q := NewSyncedQueue()
	q.Push(1)
	if q.Len() != 1 {
		t.Fatalf("Ожидалась длина %v, а получено: %v", 1, q.Len())
	}
	lastElement := q.queue[0].(int)
	if lastElement != 1 {
		t.Fatalf("Ожидалось получить %v, получено: %v", 1, lastElement)
	}
}

func TestQueueLen(t *testing.T) {
	q := getQueue()
	if q.Len() != 3 {
		t.Fatalf("Ожидалась длина %v, а получено: %v", 3, q.Len())
	}
}

func TestQueuePopFromNonEmpty(t *testing.T) {
	q := getQueue()
	element, err := q.Pop()
	if err != nil {
		t.Fatalf("pop - получена ошибка: %v", err)
	}
	firstElement := element.(int)
	if firstElement != 1 {
		t.Fatalf("Ожидалось вытянуть элемент %v, а получено: %v", 1, firstElement)
	}
	if q.Len() != 2 {
		t.Fatalf("Ожидалась длина %v, а получено: %v", 2, q.Len())
	}

	element, err = q.Pop()
	if err != nil {
		t.Fatalf("pop - получена ошибка: %v", err)
	}
	secondElement := element.(int)
	if secondElement != 2 {
		t.Fatalf("Ожидалось вытянуть элемент %v, а получено: %v", 2, secondElement)
	}
	if q.Len() != 1 {
		t.Fatalf("Ожидалась длина %v, а получено: %v", 1, q.Len())
	}

	element, err = q.Pop()
	if err != nil {
		t.Fatalf("pop - получена ошибка: %v", err)
	}
	thirdElement := element.(int)
	if thirdElement != 3 {
		t.Fatalf("Ожидалось вытянуть элемент %v, а получено: %v", 3, thirdElement)
	}
	if q.Len() != 0 {
		t.Fatalf("Ожидалась длина %v, а получено: %v", 0, q.Len())
	}
}

func TestQueuePopFromEmpty(t *testing.T) {
	q := NewSyncedQueue()
	element, err := q.Pop()
	if err == nil {
		t.Fatal("При извлечении из пустой очередь ожидалось получить ошибку, а получено nil")
	}
	if element != nil {
		t.Fatalf("При извлечении из пустой очередь ожидалось получить element == nil, а получено: %v", element)
	}
	if q.Len() != 0 {
		t.Fatalf("Ожидалась длина %v, а получено: %v", 0, q.Len())
	}
}
