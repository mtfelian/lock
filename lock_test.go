package lock

import (
	"fmt"
	"testing"
	"time"
)

func TestKeyLock(t *testing.T) {

	type runParams struct {
		key   string
		sleep int
	}
	params := []*runParams{
		{"key1", 2},
		{"key2", 4},
		{"key1", 1},
		{"key2", 3},
	}

	kl := NewKeyLock()
	sleep := 0
	start := time.Now()
	for _, p := range params {
		sleep += (p.sleep + 1)
		go func(in *runParams) {
			fmt.Printf("%v Before Block  ('%v'): %v\n", time.Since(start), in.key, in.sleep)
			kl.Block(in.key)
			fmt.Printf("%v After  Block  ('%v'): %v\n", time.Since(start), in.key, in.sleep)
			time.Sleep(time.Second * time.Duration(in.sleep))
			kl.Unblock(in.key)
			fmt.Printf("%v After  Unblock('%v'): %v\n", time.Since(start), in.key, in.sleep)
		}(p)
	}

	// Использовать sync.WaitGroup чтобы дождаться окончания работы
	// всех go-рутин не получится в данном случае.
	// можно только по таймеру, так как сам WaitGroup, зараза,
	// лочится изнутри свои mux-ом и мы получаем deadlock.
	// поэтому Sleep и никак иначе (есть другие варианты?)
	time.Sleep(time.Second * time.Duration(sleep))
}
