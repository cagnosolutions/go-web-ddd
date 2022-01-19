package user

import (
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
)

type WiredUser struct {
	*UserRepository
	*UserService
	*UserController
}

func WireUser(dao webapp.DataAccesser) *WiredUser {

	// setup and "wire" user repo
	userRepo := new(UserRepository)
	userRepo.AddDataAccesser(dao)

	// setup and "wire" user service
	userService := new(UserService)
	userService.AddRepository(userRepo)

	// setup and "wire" user controller
	userController := new(UserController)
	userController.AddService(userService)

	return &WiredUser{
		UserRepository: userRepo,
		UserService:    userService,
		UserController: userController,
	}
}
