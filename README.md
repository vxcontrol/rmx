### Recursive mutex implementation by goroutine ID

In the Golang world there is the problem that getting goroutine ID for any tasks.
But you can find [sad information](https://golang.org/doc/faq#no_goroutine_id) on the [FAQ page](https://golang.org/doc/faq) of the Golang documentation in [Rob Pike's comment](https://go.googlesource.com/go/+/992ce90f662467f04dd93b3bb565bb0414f82999%5E%21/#F0) that getting ID is imposible.
We disagree with it and think that this should be present in Golang because often needs to use shared object in goroutine safe mode.

Requires go version >= 1.4 for "Fast" mode.

This project based on [this](https://github.com/modern-go/gls) and gives possibility to use recursive mutex in two modes:
- "Generic" that provides access to goroutine ID from runtime.stack call
- "Fast" that provides the same access from reflective call to "runtime" internal object

Motivation of using both modes:
- "Generic" is common mode that will be to work for any Golang version
- "Fast" is specific mode that will be to work only limited Golang version because it's using offset to search goroutine ID in "runtime" object but that work faster about 100 times than "Generic" mode

Use cases:
- Install the package

```bash
go get github.com/vxcontrol/rmx
```

- "Generic" mode if tests were failed

```go
import "github.com/vxcontrol/rmx"
	// Some code
	rm := &rmx.Mutex{ IsGeneric: true }
		// Using inside of recursive function
		rm.Lock()
		// Goroutine safe code
		rm.Unlock()
```

- "Fast" mode if tests were passed

```go
import "github.com/vxcontrol/rmx"
	// Some code
	rm := &rmx.Mutex{ IsGeneric: false } // Option value by default
		// Using inside of recursive function
		rm.Lock()
		// Goroutine safe code
		rm.Unlock()
```

- Override default waiting time before next check mutex for highload applications
Use "-1" constant value to disable waiting (just force to run schedule goroutine process)
Use "0" constant by default value to wait a minimal delay about 1ms
Use ">1" constant value to wait a longer time in milliseconds

```go
import "github.com/vxcontrol/rmx"
	// Some code
	rm := &rmx.Mutex{ TimeWaitMS: 5 } // Wait 5ms
		// Using inside of recursive function
		rm.Lock()
		// Goroutine safe code
		rm.Unlock()
```
