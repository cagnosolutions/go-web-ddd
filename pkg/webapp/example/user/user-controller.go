package user

import (
	"errors"
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"html/template"
	"net/http"
)

// UserController implements the Controller interface
// for use with the standard library http package
type UserController struct {
	userService *UserService
	tmpls       *template.Template
}

// WithService helps satisfy the Controller interface
func (con *UserController) WithService(service webapp.Service) error {
	if service == nil {
		return errors.New("got empty service")
	}
	con.userService = service.(*UserService)
	return nil
}

// RootHandler helps satisfy the Controller interface
func (con *UserController) RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user":
		con.handleBaseRequest(w, r)
	case "/user/":
		con.handleOneUser(w, r)
	case "/user/all":
		con.handleAllUsers(w, r)
	default:
		con.handleUserError(w, r)
	}
}

func (con *UserController) handleBaseRequest(w http.ResponseWriter, r *http.Request) {
	// load the user form page
	if r.Method == http.MethodGet {
		err := con.tmpls.ExecuteTemplate(w, "user.html", nil)
		if err != nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
		}
		return
	}
	// otherwise, user has posted something
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}
		// save the user that was posted
		id, err := con.userService.AddNewUser(r)
		if err != nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}
		fmt.Fprintf(w, "successfully added user, id=%d\n", id)
		return
	}
}

func (con *UserController) handleOneUser(w http.ResponseWriter, r *http.Request) {
	// handle get user by id
	// handle update user by id
}

func (con *UserController) handleAllUsers(w http.ResponseWriter, r *http.Request) {
	// handle get request only
	// return a table of all the users
}

func (con *UserController) handleUserError(w http.ResponseWriter, r *http.Request) {
	// error handler/page for any user errors
}
