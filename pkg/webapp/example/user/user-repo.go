package user

import (
	"errors"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
)

// UserRepository implements the Repository interface
// and provides methods for the service to use
type UserRepository struct {
	userDao webapp.DAO
}

// WithDAO helps satisfy the Repository interface
func (repo *UserRepository) WithDAO(daos ...webapp.DAO) error {
	if daos == nil || len(daos) < 1 {
		return errors.New("did not receive any dao")
	}
	if len(daos) > 1 {
		return errors.New("got too many daos")
	}
	repo.userDao = daos[0]
	return nil
}

func (repo *UserRepository) AddUser(u *User) (int, error) {
	// call the dao add method
	id, err := repo.userDao.Add(u)
	// could do more processing in here if we wanted to
	return id, err
}

func (repo *UserRepository) GetUser(id int, ptr *User) error {
	// call the dao get method
	e, err := repo.userDao.Get(id)
	if err != nil {
		return err
	}
	// convert entity to user
	ptr, ok := e.(*User)
	if !ok {
		return errors.New("conversion error")
	}
	return nil
}

func (repo *UserRepository) GetAllUsers(ptr []*User) (int, error) {
	// call the dao getall method
	ee, err := repo.userDao.GetAll()
	if err != nil {
		return 0, err
	}
	// convert entities to users
	for _, e := range ee {
		u, ok := e.(*User)
		if !ok {
			return 0, errors.New("conversion error")
		}
		ptr = append(ptr, u)
	}
	return len(ptr), nil
}

func (repo *UserRepository) SetUser(u *User) error {
	// call the dao set method
	err := repo.userDao.Set(u)
	// you could do more processing in here if you wanted to
	return err
}

func (repo *UserRepository) Del(id int) error {
	// call the dao del method
	err := repo.userDao.Del(id)
	// you could do more processing in here if you wanted to
	return err
}
