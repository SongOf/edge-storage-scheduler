package model

import "sync"

type void struct{}

var member void

var EdgeSetOnline map[string]void

var lock sync.Mutex

func NewEdgeSetOnline() {
	EdgeSetOnline = make(map[string]void)
}

func Set(set string) {
	lock.Lock()
	defer lock.Unlock()
	EdgeSetOnline[set] = member
}

func Delete(set string) {
	lock.Lock()
	defer lock.Unlock()
	delete(EdgeSetOnline, set)
}

func Size() int {
	return len(EdgeSetOnline)
}

func Exists(set string) bool {
	_, exists := EdgeSetOnline[set]
	return exists
}
