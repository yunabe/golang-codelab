package goroutineid

import (
	"sync"
	"testing"
)

func TestGetGoroutineID(t *testing.T) {
	mainID := GetGoroutineID()
	ch := make(chan uint64)
	done := make(chan struct{})
	m := make(map[uint64]bool)
	go func() {
		for gid := range ch {
			if gid == mainID {
				t.Errorf("Goroutine ID of a goroutine is same as the ID of the main routine: %d", gid)
			}
			if m[gid] {
				t.Errorf("Goroutine ID dup: %d", gid)
			}
			m[gid] = true
		}
		close(done)
	}()

	for i := 0; i < 1000; i++ {
		rep := 100
		var wg sync.WaitGroup
		wg.Add(rep)
		for j := 0; j < rep; j++ {
			go func() {
				ch <- GetGoroutineID()
				wg.Done()
			}()
		}
		wg.Wait()
	}
	close(ch)
	<-done

	if len(m) != 100000 {
		t.Errorf("Unexpected len(m): %d", len(m))
		return
	}
}
