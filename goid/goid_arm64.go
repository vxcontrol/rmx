package goid

func getg() *g

func Get() int64 {
	return getg().goid
}
