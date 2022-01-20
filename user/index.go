package user

import (
	"github.com/chuccp/cokePush/util"
	"sync"
	"sync/atomic"
	"time"
)

type StoreUser struct {
	store *sync.Map
	num   int32
	username string
	createTime *time.Time
}

func newUserStore(username string) *StoreUser {
	t:=time.Now()
	return &StoreUser{store: new(sync.Map),username:username,createTime:&t}
}
func (userStore *StoreUser) add(user IUser) (bool, int32) {
	_, fa := userStore.store.LoadOrStore(user.GetId(), user)
	if !fa {
		return true, atomic.AddInt32(&userStore.num, 1)
	}
	return false, userStore.num
}
func (userStore *StoreUser) delete(user IUser) (bool, int32) {
	id := user.GetId()
	_, fa := userStore.store.LoadAndDelete(id)
	if fa {
		return true, atomic.AddInt32(&userStore.num, -1)
	}
	return false, userStore.num
}
func (userStore *StoreUser) each(f func(IUser) bool) {
	userStore.store.Range(func(key, value interface{}) bool {
		return f(value.(IUser))
	})

}

func (userStore *StoreUser) GetUsername() string {
	return userStore.username
}

func (userStore *StoreUser) MachineAddress() string {
	return "localhost"
}

func (userStore *StoreUser) CreateTime() string {
	return userStore.createTime.Format(util.TimestampFormat)
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
		us := v.(*StoreUser)
		us.add(user)
	} else {
		us := newUserStore(username)
		us.add(user)
		index.uMap.Store(username, us)
		atomic.AddInt32(&index.num, 1)
		index.rLock.Unlock()
		return true
	}
	index.rLock.Unlock()
	return false
}
func (index *indexMap) eachUsers(f func(username string,user *StoreUser)bool){
	index.uMap.Range(func(key, value interface{}) bool {
		return f(key.(string),value.(*StoreUser))
	})
}
func (index *indexMap) delete(user IUser) bool {
	username := user.GetUsername()
	v, ok := index.uMap.Load(username)
	if ok {
		us := v.(*StoreUser)
		_, num := us.delete(user)
		if num == 0 {
			index.rLock.Lock()
			if us.num == 0 {
				index.uMap.Delete(username)
				atomic.AddInt32(&index.num, -1)
			}
			index.rLock.Unlock()
			return true
		}

	}
	return false
}
func (index *indexMap)has(username string)bool{
	 _,ok:=index.uMap.Load(username)
	 return ok
}
func (index *indexMap) each(username string, f func(IUser) bool) bool {
	v, ok := index.uMap.Load(username)
	if ok {
		us, ok := v.(*StoreUser)
		if ok {
			us.each(f)
		}
	}
	return ok
}
