package user

import "sync"

type indexMap struct {
	uMap  *sync.Map
	num   int
	rLock *sync.RWMutex
}

func newIndexMap() *indexMap {
	return &indexMap{new(sync.Map), 0, new(sync.RWMutex)}
}

func (index *indexMap) add(user IUser) bool {
	username := user.GetUsername()
	index.rLock.Lock()
	index.uMap.Store(username, user)
	index.rLock.Unlock()
	return false
}
func (index *indexMap) get(username string) IUser {
	index.rLock.RLock()
	u, ok := index.uMap.Load(username)
	index.rLock.RUnlock()
	if ok {
		return u.(IUser)
	} else {
		return nil
	}
}
