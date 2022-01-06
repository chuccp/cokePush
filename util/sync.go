package util

import "sync"

type MapLock struct {
	lMap  map[string]*sync.RWMutex
	rLock *sync.RWMutex
}
func NewMapLock() *MapLock {
	return &MapLock{lMap: make(map[string]*sync.RWMutex), rLock: new(sync.RWMutex)}
}
func (m *MapLock) RLock(key string) {
	m.rLock.RLock()
	sr := m.lMap[key]
	if sr == nil {
		m.rLock.RUnlock()
		m.rLock.Lock()
		sr := m.lMap[key]
		if sr == nil {
			rm := new(sync.RWMutex)
			m.lMap[key] = rm
			m.rLock.Unlock()
			rm.RLock()
		} else {
			m.rLock.Unlock()
			sr.RLock()
		}
	} else {
		m.rLock.RUnlock()
		sr.RLock()
	}
}
func (m *MapLock) Lock(key string) {
	m.rLock.RLock()
	sr := m.lMap[key]
	if sr == nil {
		m.rLock.RUnlock()
		m.rLock.Lock()
		sr := m.lMap[key]
		if sr == nil {
			rm := new(sync.RWMutex)
			m.lMap[key] = rm
			m.rLock.Unlock()
			rm.Lock()
		} else {
			m.rLock.Unlock()
			sr.Lock()
		}
	} else {
		m.rLock.RUnlock()
		sr.Lock()
	}

}

func (m *MapLock) RUnLock(key string) {
	m.rLock.RLock()
	sr := m.lMap[key]
	if sr == nil {
		m.rLock.RUnlock()
	} else {
		m.rLock.RUnlock()
		m.rLock.Lock()
		sr := m.lMap[key]
		if sr == nil {
			m.rLock.Unlock()
		} else {
			m.rLock.Unlock()
			sr.RUnlock()
		}
	}
}
func (m *MapLock) UnLock(key string) {
	m.rLock.RLock()
	sr := m.lMap[key]
	if sr == nil {
		m.rLock.RUnlock()
	} else {
		m.rLock.RUnlock()
		m.rLock.Lock()
		sr := m.lMap[key]
		if sr == nil {
			m.rLock.Unlock()
		} else {
			m.rLock.Unlock()
			sr.Unlock()
		}
	}
}
