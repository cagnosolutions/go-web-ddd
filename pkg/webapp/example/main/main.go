package main

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/example/user"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/memory"
	"net/http"
)

func main() {

	// init dao (data source)
	dao := memory.NewMemoryDataSource()
	loadUser(dao)

	// test path
	paths := []string{"/", "/user", "/user/", "/user/1", "/user/1/"}
	for i := range paths {
		p := webapp.NewPath(paths[i])
		fmt.Printf("original=%q\t\tpath=%q\tid=%q\thasID=%v\n",
			paths[i], p.Path, p.ID, p.HasID())
	}

}

func loadUser(dao webapp.DataAccesser) (string, http.Handler) {
	// init and wire user
	userController := user.WireUser(dao)

	// add to your main router wherever that is
	return "/user", http.HandlerFunc(userController.HandleBase)
}
