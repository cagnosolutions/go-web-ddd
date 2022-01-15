package user

import (
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
)

func WireUser(dao webapp.DAO) *UserController {

	// setup and "wire" user repo
	userRepo := new(UserRepository)
	err := userRepo.WithDAO(dao)
	if err != nil {
		panic(err)
	}

	// setup and "wire" user service
	userService := new(UserService)
	err = userService.WithRepo(userRepo)
	if err != nil {
		panic(err)
	}

	// setup and "wire" user controller
	userController := new(UserController)
	err = userController.WithService(userService)
	if err != nil {
		panic(err)
	}

	return userController
}
