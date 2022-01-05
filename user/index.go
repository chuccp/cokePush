package user

import (
	"sync"
)

type userStore struct {
	store *sync.Map
	num   int
}

func newUserStore() *userStore {
	return &userStore{store: new(sync.Map)}
}
func (userStore *userStore) add(user IUser) (bool,int) {
	_, fa := userStore.store.LoadOrStore(user.GetUserId(), user)
	if !fa {
		userStore.num++
		return true,userStore.num
	}
	return false,userStore.num
}
func (userStore *userStore) delete(user IUser) (bool,int) {
	id := user.GetUserId()
	_, fa := userStore.store.LoadAndDelete(id)
	if fa {
		userStore.num--
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
	num   uint
	rLock *sync.RWMutex
}

func newIndexMap() *indexMap {
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
		index.num++
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
				index.num--
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