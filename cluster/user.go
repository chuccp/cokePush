package cluster

import (
	"github.com/chuccp/cokePush/message"
	"strconv"
	"sync"
)
type cu struct {
	userId string
	username string
	remoteHost string
	remotePort int
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
func(u *cu)SetUsername(username string){
	u.username = username
	u.userId = username+u.remoteHost+strconv.Itoa(u.remotePort)
}
func newCu(username string,remoteHost string,remotePort int)*cu  {
	u:=&cu{username:username,remoteHost:remoteHost,remotePort:remotePort}
	u.SetUsername(username)
	return u
}

type userStore struct {
	userMap * sync.Map
}

func (us *userStore)AddUser(username string,remoteHost string,remotePort int)  {
	c:=newCu(username,remoteHost,remotePort)
	us.userMap.Store(username,c)
}
func (us *userStore)DeleteUser(username string)  {
	us.userMap.Delete(username)
}
func (us *userStore)GetUserMachine(username string)(*cu,bool)  {
	v,ok:=us.userMap.Load(username)
	if ok{
		return v.(*cu),ok
	}
	return nil,ok
}
func newUserStore() *userStore {
	return &userStore {userMap: new(sync.Map)}
}