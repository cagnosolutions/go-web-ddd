package user

import (
	"errors"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"net/http"
)

// UserService implements the Service interface
// and provides methods for the controller to use
type UserService struct {
	userRepo *UserRepository
}

// WithRepo helps satisfy the Service interface
func (service *UserService) WithRepo(repo webapp.Repository) error {
	if repo == nil {
		return errors.New("got empty repo")
	}
	service.userRepo = repo.(*UserRepository)
	return nil
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
