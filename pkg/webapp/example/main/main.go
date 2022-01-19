package main

import (
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/example/user"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/memory"
	"log"
	"net/http"
)

var (
	tc *webapp.TemplateCache
	ss *webapp.SessionStore
	db webapp.DataAccesser
)

func init() {
	// init templates
	tc = webapp.NewTemplateCache("pkg/webapp/example/main/web/templates/*.html", nil)
	tc.ParseGlob("pkg/webapp/example/main/web/templates/stubs/*.html")

	// init session store
	ss = webapp.NewSessionStore("sess-id", 300)

	// init in memory db
	db = memory.NewMemoryDataSource()
}

func main() {

	// server
	mux := http.NewServeMux()
	mux.Handle("/error/", webapp.ErrorHandler(tc.Lookup("error.html")))
	mux.Handle("/index", handleIndex(tc))
	mux.Handle("/login", handleLogin(tc, ss))
	mux.Handle("/logout", handleLogout(ss))
	mux.Handle("/secure/home", handleSecureHome(ss))
	mux.Handle("/templates", handleTemplates(tc))
	mux.Handle("/bootstrap", handleBootstrapExample())
	mux.Handle("/static/", webapp.StaticHandler("/static", "pkg/webapp/example/main/web/static/"))
	log.Fatal(http.ListenAndServe(":8080", mux))

}

func loadUser(dao webapp.DataAccesser) (string, http.Handler) {
	// init and wire user
	userController := user.WireUser(dao)

	// add to your main router wherever that is
	return "/user", userController.HandleBase()
}
