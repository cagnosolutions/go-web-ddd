package domain

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/user-project/resources/file/db"
)

type UserRepository interface {
	Get(id int) (*User, error)
	GetAll() ([]*User, error)
}

func NewUserRepository() *DefaultUserRepository {
	db, err := db.OpenFileDB("mydata/user-data.csv")
	if err != nil {
		panic(fmt.Sprintf("user-repo: %s", err))
	}
	return &DefaultUserRepository{db: db}
}

type DefaultUserRepository struct {
	db *db.FileDB
}

func (r *DefaultUserRepository) Get(id int) (*User, error) {
	u, err := r.db.Get(id)
	if err != nil {
		return nil, err
	}
	return u.(*User), nil
}

func (r *DefaultUserRepository) GetAll() ([]*User, error) {
	us, err := r.db.GetAll()
	if err != nil {
		return nil, err
	}
	var u []*User
	for i := range us {
		u = append(u, us[i].(*User))
	}
	return u, nil
}
