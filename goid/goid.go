package goid

import (
	"bytes"
	"runtime"
	"strconv"
)

func SlowGet() int64 {
	var buf [32]byte
	l := runtime.Stack(buf[:], false)             // fill buffer from stack trace info
	s := buf[10:l]                                // trim "goroutine" prefix
	s = s[:bytes.IndexByte(s, ' ')]               // trim everything after first space
	gid, _ := strconv.ParseInt(string(s), 10, 64) // parsae as int value
	return gid
}
