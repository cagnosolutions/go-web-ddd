package controller

import (
	"encoding/json"
	"github.com/cagnosolutions/go-web-ddd/user-project/service"
	"net/http"
)

type UserHandler struct {
	Service service.UserService
}

func (handler *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := handler.Service.GetAllUsers()

	// setting header type is must otherwise response will be plain/text format
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
