package rmx

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/modern-go/gls"
)

// Mutex is struct that contains information about recursive calls and input args:
// IsGeneric is initial argument that set behaviour of receiving goroutine ID
// TimeWaitMS is initial argument that set waiting time between mutex state checks
type Mutex struct {
	IsGeneric        bool
	TimeWaitMS       time.Duration
	mutex            sync.Mutex
	internalMutex    sync.Mutex
	currentGoRoutine int64
	lockCount        uint64
}

func (m *Mutex) getGenericRoutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	if n <= 0 {
		return 0
	}
	idFields := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))
	if len(idFields) <= 0 {
		return 0
	}
	id, err := strconv.ParseInt(idFields[0], 10, 64)
	if err != nil {
		return 0
	}

	return id
}

func (m *Mutex) getFastRoutineID() int64 {
	return gls.GoID()
}

func (m *Mutex) getRoutineID() int64 {
	if m.IsGeneric {
		return m.getGenericRoutineID()
	}

	return m.getFastRoutineID()
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

// Lock is function that locks state of current goroutine
func (m *Mutex) Lock() {
	goRoutineID := m.getRoutineID()

	for {
		m.internalMutex.Lock()
		if m.currentGoRoutine == 0 {
			m.currentGoRoutine = goRoutineID
			break
		} else if m.currentGoRoutine == goRoutineID {
			break
		} else {
			m.internalMutex.Unlock()
			runtime.Gosched()
			m.wait()
			continue
		}
	}
	m.lockCount++
	m.internalMutex.Unlock()
}

// Unlock is function that unlocks state of current goroutine
func (m *Mutex) Unlock() {
	goRoutineID := m.getRoutineID()

	for {
		m.internalMutex.Lock()
		if m.currentGoRoutine == 0 {
			m.currentGoRoutine = goRoutineID
			break
		} else if m.currentGoRoutine == goRoutineID {
			break
		} else {
			m.internalMutex.Unlock()
			runtime.Gosched()
			m.wait()
			continue
		}
	}
	m.lockCount--
	if m.lockCount == 0 {
		m.currentGoRoutine = 0
	}
	m.internalMutex.Unlock()
}
