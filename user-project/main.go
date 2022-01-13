package main

import (
	"github.com/cagnosolutions/go-web-ddd/user-project/controller"
	"github.com/cagnosolutions/go-web-ddd/user-project/domain"
	"github.com/cagnosolutions/go-web-ddd/user-project/service"
	"log"
	"net/http"
)

func main() {

	// Define routes
	mux := http.NewServeMux()

	// wiring
	// note that we have decided to use stub repo and default service by
	// instantiating them.
	ch := controller.UserHandler{Service: service.NewUserService(domain.NewUserRepository())}
	mux.HandleFunc("/users", ch.GetAllUsers)

	// Starting the server
	log.Fatal(http.ListenAndServe("localhost:8001", mux))
}
