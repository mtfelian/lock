package lock

import (
	"fmt"
	"testing"
	"time"
)

// TestKeyLock checks keyLock
func TestKeyLock(t *testing.T) {
	type testCases struct {
		key   string
		sleep int
	}
	params := []*testCases{
		{"key1", 2},
		{"key2", 4},
		{"key1", 1},
		{"key2", 3},
		{"key1", 2},
		{"key2", 3},
		{"key1", 2},
		{"key2", 1},
	}

	kl := NewKeyLock()
	sleep := 0
	start := time.Now()
	for _, p := range params {
		sleep += (p.sleep + 1)
		go func(in *testCases) {
			fmt.Printf("%v Before Block  ('%v'): %v\n", time.Since(start), in.key, in.sleep)
			kl.Block(in.key)
			fmt.Printf("%v After  Block  ('%v'): %v\n", time.Since(start), in.key, in.sleep)
			time.Sleep(time.Second * time.Duration(in.sleep) / 3)
			kl.Unblock(in.key)
			fmt.Printf("%v After  Unblock('%v'): %v\n", time.Since(start), in.key, in.sleep)
		}(p)
	}

	// We can not use sync.WaitGroup here to wait for all goroutines
	// finished due to WaitGroup uses mutex internally and this goes into a deadlock.
	// So we use time.Sleep(). Altogether exists a better solution here
	// todo remove time.Sleep(), rework this test
	time.Sleep(time.Second * time.Duration(sleep) / 3)
}

// TestKeyLockUnblockTwice tests twice unblocking of keyLock without panic
func TestKeyLockUnblockTwice(t *testing.T) {
	keyLock := NewKeyLock()
	key := "key"
	keyLock.Block(key)
	keyLock.Unblock(key)
	keyLock.Unblock(key)
}
