package rmx

import (
	"fmt"
	"sync"
	"testing"
)

func recfunc(rm *Mutex, syncInt *int, n int) {
	var finish bool

	for {
		rm.Lock()
		switch *syncInt % 4 {
		case 0:
			*syncInt++
			recfunc(rm, syncInt, 0)
		case 2:
			*syncInt++
			recfunc(rm, syncInt, 0)
		case n:
			*syncInt++
			finish = true
		default:
		}
		rm.Unlock()

		if n == 0 || finish {
			break
		}
	}
}

func TestNewMutexDefault(t *testing.T) {
	rm := &Mutex{}
	rm.Lock()
	_ = rm.Unlock()
}

func TestNewMutexGenericMode(t *testing.T) {
	rm := &Mutex{IsGeneric: true}
	rm.Lock()
	_ = rm.Unlock()
}

func TestNewMutexFastMode(t *testing.T) {
	rm := &Mutex{IsGeneric: false}
	rm.Lock()
	_ = rm.Unlock()
}

func TestGettingGoroutineID(t *testing.T) {
	rm1 := &Mutex{IsGeneric: true}
	rm2 := &Mutex{IsGeneric: false}
	var wg sync.WaitGroup
	failed := false
	runner := func() {
		defer wg.Done()
		if rm1.getRoutineID() != rm2.getRoutineID() {
			t.Log("Failed retrieve goroutine ID correct")
			failed = true
		}
	}

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go runner()
	}
	wg.Wait()
	if failed {
		t.FailNow()
	}
}

// Run simple test with shared object
func ExampleMutex_default() {
	rm := &Mutex{}
	var syncInt int
	var wg sync.WaitGroup
	runner := func(n int) {
		defer wg.Done()
		recfunc(rm, &syncInt, n)
	}

	for i := 0; i < 500; i++ {
		wg.Add(2)
		go runner(1)
		go runner(3)
		wg.Wait()
	}
	fmt.Println(syncInt)
	// Output:
	// 2000
}

// Run simple test with shared object in generic mode
func ExampleMutex_generic_mode() {
	rm := &Mutex{IsGeneric: true}
	var syncInt int
	var wg sync.WaitGroup
	runner := func(n int) {
		defer wg.Done()
		recfunc(rm, &syncInt, n)
	}

	for i := 0; i < 500; i++ {
		wg.Add(2)
		go runner(1)
		go runner(3)
		wg.Wait()
	}
	fmt.Println(syncInt)
	// Output:
	// 2000
}

// Run simple test with shared object in fast mode
func ExampleMutex_fast_mode() {
	rm := &Mutex{IsGeneric: false}
	var syncInt int
	var wg sync.WaitGroup
	runner := func(n int) {
		defer wg.Done()
		recfunc(rm, &syncInt, n)
	}

	for i := 0; i < 500; i++ {
		wg.Add(2)
		go runner(1)
		go runner(3)
		wg.Wait()
	}
	fmt.Println(syncInt)
	// Output:
	// 2000
}

func BenchmarkMutex_getRoutineID_generic_mode(b *testing.B) {
	rm := &Mutex{IsGeneric: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.getRoutineID()
	}
}

func BenchmarkMutex_getRoutineID_fast_mode(b *testing.B) {
	rm := &Mutex{IsGeneric: false}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.getRoutineID()
	}
}

func BenchmarkMutex_generic_mode(b *testing.B) {
	rm := &Mutex{IsGeneric: true}
	var syncInt int
	var wg sync.WaitGroup
	runner := func(n int) {
		defer wg.Done()
		recfunc(rm, &syncInt, n)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go runner(1)
		go runner(3)
		wg.Wait()
	}
}

func BenchmarkMutex_fast_mode(b *testing.B) {
	rm := &Mutex{IsGeneric: false}
	var syncInt int
	var wg sync.WaitGroup
	runner := func(n int) {
		defer wg.Done()
		recfunc(rm, &syncInt, n)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go runner(1)
		go runner(3)
		wg.Wait()
	}
}

func TestLockFromOtherGoroutine(t *testing.T) {
	m := Mutex{}
	s := make(chan struct{})

	wg := sync.WaitGroup{}
	failed := false

	wg.Add(2)
	go func() {
		defer wg.Done()
		m.Lock()
		s <- struct{}{}
		if err := m.Unlock(); err != nil {
			t.Logf("unexpected error: %v", err)
			failed = true
		}
	}()

	go func() {
		defer wg.Done()
		<-s
		m.Lock()
		if err := m.Unlock(); err != nil {
			t.Logf("unexpected error: %v", err)
			failed = true
		}
	}()

	wg.Wait()

	if failed {
		t.FailNow()
	}
}

func TestUnlockFail(t *testing.T) {
	m := Mutex{}
	s := make(chan struct{})

	wg := sync.WaitGroup{}
	failed := false

	wg.Add(2)
	go func() {
		defer wg.Done()
		m.Lock()
		if err := m.Unlock(); err != nil {
			t.Logf("unexpected error: %v", err)
			failed = true
		}
		s <- struct{}{}
	}()

	go func() {
		defer wg.Done()
		<-s
		if err := m.Unlock(); err != ErrAlreadyUnlocked {
			t.Logf("error expected")
			failed = true
		}
	}()

	wg.Wait()

	if failed {
		t.FailNow()
	}
}

func TestOtherGoroutineLock(t *testing.T) {
	m := Mutex{}
	s := make(chan struct{})

	wg := sync.WaitGroup{}
	failed := false

	wg.Add(2)
	go func() {
		defer wg.Done()
		m.Lock()
		s <- struct{}{}
		<-s
		_ = m.Unlock()
	}()

	go func() {
		defer wg.Done()
		<-s
		if err := m.Unlock(); err != ErrOtherGoroutineLocked {
			t.Logf("error expected")
			failed = true
		}
		s <- struct{}{}
	}()

	wg.Wait()

	if failed {
		t.FailNow()
	}
}
