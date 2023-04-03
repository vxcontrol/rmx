package rmx

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/vxcontrol/rmx/goid"
)

// Mutex goroutine-bound mutex implementation
type Mutex struct {
	// IsGeneric determines if slow (but probably more guaranteed way) goroutine ID detection
	// process should be used
	IsGeneric bool
	// TimeWaitMS defines period between mutex state checks
	TimeWaitMS       time.Duration
	mutex            sync.Mutex
	currentGoRoutine int64
	lockCount        uint64
}

var (
	ErrAlreadyUnlocked      = errors.New("Goroutine already unlocked")
	ErrOtherGoroutineLocked = errors.New("Other goroutine already take the lock")
)

func (m *Mutex) getRoutineID() int64 {
	if m.IsGeneric {
		return goid.SlowGet()
	}

	return goid.Get()
}

func (m *Mutex) wait() {
	switch twait := m.TimeWaitMS; {
	case twait == 0:
		time.Sleep(time.Millisecond)
	case twait > 0:
		time.Sleep(time.Millisecond * m.TimeWaitMS)
	default:
	}
}

// Lock locks state of the current goroutine
func (m *Mutex) Lock() {
	goRoutineID := m.getRoutineID()

	for {
		m.mutex.Lock()
		if m.currentGoRoutine == 0 {
			m.currentGoRoutine = goRoutineID
			break
		}
		if m.currentGoRoutine == goRoutineID {
			break
		}
		m.mutex.Unlock()

		runtime.Gosched()
		m.wait()
	}
	m.lockCount++
	m.mutex.Unlock()
}

// Unlock unlocks state of the current goroutine
func (m *Mutex) Unlock() error {
	goRoutineID := m.getRoutineID()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.currentGoRoutine == 0 {
		return ErrAlreadyUnlocked
	}
	if m.currentGoRoutine != goRoutineID {
		return ErrOtherGoroutineLocked
	}

	m.lockCount--
	if m.lockCount == 0 {
		m.currentGoRoutine = 0
	}

	return nil
}
