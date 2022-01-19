package user

import (
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"net/http"
)

// UserService implements the Servicer interface
// and provides methods for the controller to use
type UserService struct {
	userRepo *UserRepository
}

// AddRepository helps satisfy the Servicer interface
func (service *UserService) AddRepository(repo webapp.Repository) {
	if repo == nil {
		panic("got empty repo")
	}
	service.userRepo = repo.(*UserRepository)
}

func (service *UserService) AddNewUser(r *http.Request) (int, error) {
	// get the posted user
	user := &User{
		FirstName:    r.FormValue("first"),
		LastName:     r.FormValue("last"),
		EmailAddress: r.FormValue("email"),
	}
	user.UpdatePassword(r.FormValue("password"))
	// save the new user to the database
	return service.userRepo.AddUser(user)
}

func (service *UserService) GetUser(un, pw string) *User {
	var users []*User
	count, err := service.userRepo.GetAllUsers(users)
	if err != nil || count < 1 {
		return nil
	}
	for i := range users {
		if users[i].EmailAddress == un && users[i].Password == pw {
			return users[i]
		}
	}
	return nil
}
