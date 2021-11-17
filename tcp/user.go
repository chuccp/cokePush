package tcp

import (
	"github.com/chuccp/cokePush/user"
)
type User struct {
	*user.User
}

func NewUser() *User {
	return &User{}
}