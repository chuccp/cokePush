package cluster

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)
type cu struct {
	userId string
	username string
	remoteHost string
	remotePort int
	machineAddress string
	createTime *time.Time
}
func(u *cu)WriteMessage(iMessage message.IMessage) error{
	return nil
}
func(u *cu)GetId() string{
	return u.userId
}
func(u *cu)GetUsername() string{
	return u.username
}
func(u *cu)MachineAddress() string{
	return u.machineAddress
}
func(u *cu)CreateTime() string{
	return u.createTime.Format(util.TimestampFormat)
}
func(u *cu)SetUsername(username string){
	u.username = username
	u.userId = username+u.remoteHost+strconv.Itoa(u.remotePort)
}
func newCu(username string,remoteHost string,remotePort int)*cu  {
	u:=&cu{username:username,remoteHost:remoteHost,remotePort:remotePort}
	u.SetUsername(username)
	u.machineAddress = remoteHost+":"+strconv.Itoa(remotePort)
	tu:=time.Now()
	u.createTime = &tu
	return u
}

type userStore struct {
	userMap * sync.Map
	num int32
}

func (us *userStore)AddUser(username string,remoteHost string,remotePort int)  {
	c:=newCu(username,remoteHost,remotePort)
	us.userMap.Store(username,c)
	atomic.AddInt32(&us.num,1)
}
func (us *userStore)DeleteUser(username string)  {
	us.userMap.Delete(username)
	atomic.AddInt32(&us.num,-1)
}
func (us *userStore)Num()int32  {
	return us.num
}
func (us *userStore)GetUserMachine(username string)(*cu,bool)  {
	v,ok:=us.userMap.Load(username)
	if ok{
		return v.(*cu),ok
	}
	return nil,ok
}
func (us *userStore)EachUsers(f func(key string, value *cu) bool)  {
	us.userMap.Range(func(key, value interface{}) bool {
		return f(key.(string),value.(*cu))
	})
}

func newUserStore() *userStore {
	return &userStore {userMap: new(sync.Map)}
}