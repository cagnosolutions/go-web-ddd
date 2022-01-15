package user

import (
	"crypto/sha256"
)

// User is a user model
type User struct {
	ID           int
	FirstName    string
	LastName     string
	EmailAddress string
	Password     [32]byte
	IsActive     bool
}

func NewUser(fname, lname, email string) *User {
	return &User{
		FirstName:    fname,
		LastName:     lname,
		EmailAddress: email,
		IsActive:     true,
	}
}

func (u *User) UpdatePassword(pass string) {
	u.Password = sha256.Sum256([]byte(pass))
}

// GetID helps satisfy the Entity interface
func (u *User) GetID() int {
	return u.ID
}

// SetID helps satisfy the Entity interface
func (u *User) SetID(id int) {
	if u.ID != id {
		u.ID = id
	}
}
