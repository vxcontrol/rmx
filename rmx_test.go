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
	rm.Unlock()
}

func TestNewMutexGenericMode(t *testing.T) {
	rm := &Mutex{IsGeneric: true}
	rm.Lock()
	rm.Unlock()
}

func TestNewMutexFastMode(t *testing.T) {
	rm := &Mutex{IsGeneric: false}
	rm.Lock()
	rm.Unlock()
}

func TestGettingGoroutineID(t *testing.T) {
	rm1 := &Mutex{IsGeneric: true}
	rm2 := &Mutex{IsGeneric: false}
	var wg sync.WaitGroup
	runner := func() {
		defer wg.Done()
		if rm1.getRoutineID() != rm2.getRoutineID() {
			t.Fatal("Failed retrieve goroutine ID correct")
		}
	}

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go runner()
	}
	wg.Wait()
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
