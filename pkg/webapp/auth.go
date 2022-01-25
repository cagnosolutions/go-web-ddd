package webapp

import "sync"

type AuthUser interface {
	Register(username, password, role string)
	Authenticate(username, password string) (*SystemUser, bool)
}

type SystemUser struct {
	Username string
	Password string
	Role     string
}

type SystemSessionUser struct {
	users *sync.Map
}

func NewSystemSessionUser() *SystemSessionUser {
	return &SystemSessionUser{
		users: new(sync.Map),
	}
}

func (a *SystemSessionUser) Register(username, password, role string) {
	a.users.Store(username, &SystemUser{
		Username: username,
		Password: password,
		Role:     role,
	})
}

func (a *SystemSessionUser) Authenticate(username, password string) (*SystemUser, bool) {
	su, ok := a.users.Load(username)
	if !ok {
		return nil, false
	}
	if su.(*SystemUser).Password != password {
		return nil, false
	}
	return su.(*SystemUser), true
}
