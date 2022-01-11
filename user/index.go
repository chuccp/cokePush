package user

import (
	"sync"
	"sync/atomic"
)

type userStore struct {
	store *sync.Map
	num   int32
}

func newUserStore() *userStore {
	return &userStore{store: new(sync.Map)}
}
func (userStore *userStore) add(user IUser) (bool,int32) {
	_, fa := userStore.store.LoadOrStore(user.GetId(), user)
	if !fa {
		atomic.AddInt32(&userStore.num,1)
		return true,userStore.num
	}
	return false,userStore.num
}
func (userStore *userStore) delete(user IUser) (bool,int32) {
	id := user.GetId()
	_, fa := userStore.store.LoadAndDelete(id)
	if fa {
		atomic.AddInt32(&userStore.num,-1)
		return true,userStore.num
	}
	return false,userStore.num
}
func (userStore *userStore)each(f func(IUser)bool)  {
	userStore.store.Range(func(key, value interface{}) bool {
		return f(value.(IUser))
	})

}

type indexMap struct {
	uMap  *sync.Map
	num   int32
	rLock *sync.RWMutex
}

func NewIndexMap() *indexMap {
	return &indexMap{new(sync.Map), 0, new(sync.RWMutex)}
}

func (index *indexMap) add(user IUser) bool {
	username := user.GetUsername()
	index.rLock.Lock()
	v, ok := index.uMap.Load(username)
	if ok {
		us, ok := v.(*userStore)
		if ok {
			us.add(user)
		}
	}else{
	 	us:=newUserStore()
		us.add(user)
		index.uMap.Store(username, us)
		 atomic.AddInt32(&index.num,1)
		index.rLock.Unlock()
		 return true
	}
	index.rLock.Unlock()
	return false
}

func (index *indexMap) delete(user IUser) bool {
	username := user.GetUsername()
	index.rLock.Lock()
	v, ok := index.uMap.Load(username)
	if ok {
		us, ok := v.(*userStore)
		if ok {
			_,num:=us.delete(user)
			if num==0{
				index.uMap.Delete(username)
				atomic.AddInt32(&index.num,-1)
				index.rLock.Unlock()
				return true
			}
		}
	}
	index.rLock.Unlock()
	return false
}
func (index *indexMap)each(username string,f func(IUser)bool)bool{
	v, ok := index.uMap.Load(username)
	if ok{
		us, ok := v.(*userStore)
		if ok{
			us.each(f)
		}
	}
	return ok
}