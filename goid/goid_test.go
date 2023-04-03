package goid

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSlowGet(t *testing.T) {
	if SlowGet() == 0 {
		t.FailNow()
	}
}

func TestGet(t *testing.T) {
	const count = 10000
	wg := sync.WaitGroup{}

	errs := make(chan error, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Int31n(100)))

			actual := Get()
			expected := SlowGet()
			if actual != expected {
				errs <- fmt.Errorf("Expected %d, but got %d", expected, actual)
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatal(err)
	}
}
