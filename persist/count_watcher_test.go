package persist

import (
	"testing"
)

func TestCountWatcher(t *testing.T) {
	watcher := NewCountWatcher()

	c := make(chan interface{})
	go func() {
		for i := 0; i < 20; i++ {
			watcher.Login("", "", "", nil)
		}

		c <- struct{}{}
	}()

	for i := 0; i < 18; i++ {
		watcher.Logout("", "", "")
	}
	<-c
	t.Logf("GetCounts = %v", watcher.GetLoginCounts())
}
