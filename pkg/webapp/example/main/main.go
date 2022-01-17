package main

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/example/user"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/memory"
	"log"
	"net/http"
)

func main() {

	// init template
	tc := webapp.NewTemplateCache("pkg/webapp/example/main/templates/*.html", nil)
	fmt.Println(tc.DefinedTemplates())
	tc.ParseGlob("pkg/webapp/example/main/templates/stubs/*.html")
	fmt.Println(tc.DefinedTemplates())

	// server
	mux := http.NewServeMux()
	mux.Handle("/index", handleIndex(tc))
	mux.Handle("/login", handleLogin(tc))
	mux.Handle("/templates", handleTemplates(tc))
	log.Fatal(http.ListenAndServe(":8080", mux))

}

func handleIndex(t *webapp.TemplateCache) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.html", map[string]interface{}{})
	}
	return http.HandlerFunc(fn)
}

func handleLogin(t *webapp.TemplateCache) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "login.html", map[string]interface{}{})
	}
	return http.HandlerFunc(fn)
}

func handleTemplates(t *webapp.TemplateCache) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", t.DefinedTemplates())
		return
	}
	return http.HandlerFunc(fn)
}

func testing() {
	// init dao (data source)
	dao := memory.NewMemoryDataSource()
	loadUser(dao)

	// test path
	paths := []string{"/", "/user", "/user/", "/user/1", "/user/1/", "/user/3/order/24"}
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
	return "/user", userController.HandleBase()
}
