package user

import (
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
)

func WireUser(dao webapp.DAO) *UserController {

	// setup and "wire" user repo
	userRepo := new(UserRepository)
	userRepo.AddDAO(dao)

	// setup and "wire" user service
	userService := new(UserService)
	userService.AddRepository(userRepo)

	// setup and "wire" user controller
	userController := new(UserController)
	userController.AddService(userService)

	return userController
}
