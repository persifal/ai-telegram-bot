package telegram

import (
	"strconv"
	"sync"
)

type whitelist struct {
	sync.RWMutex
	users map[int64]bool
}

func newWhiteList(names []string) *whitelist {
	m := make(map[int64]bool)
	for _, v := range names {
		i, _ := strconv.Atoi(v)
		m[int64(i)] = true
	}

	whitelist := &whitelist{users: m}
	return whitelist
}

func (a *whitelist) contains(id int64) bool {
	a.RLock()
	defer a.RUnlock()
	return a.users[id]
}
