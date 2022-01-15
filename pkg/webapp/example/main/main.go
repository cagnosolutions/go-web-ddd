package main

import (
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/example/user"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/memory"
	"net/http"
)

func main() {

	// init dao (data source)
	dao := memory.NewMemoryDataSource()

	// init and wire user
	userController := user.WireUser(dao)

	// add to your main router wherever that is
	http.HandleFunc("/user", userController.RootHandler)
}
