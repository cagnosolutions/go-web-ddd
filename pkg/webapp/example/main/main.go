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
	tc := webapp.NewTemplateCache("pkg/webapp/example/main/web/templates/*.html", nil)
	fmt.Println(tc.DefinedTemplates())
	tc.ParseGlob("pkg/webapp/example/main/web/templates/stubs/*.html")
	fmt.Println(tc.DefinedTemplates())

	// init session store
	ss := webapp.NewSessionStore("sess-id", 60)

	// server
	mux := http.NewServeMux()
	mux.Handle("/error/", webapp.ErrorHandler(tc.Lookup("error.html")))
	mux.Handle("/index", handleIndex(tc))
	mux.Handle("/login", handleLogin(tc, ss))
	mux.Handle("/logout", handleLogout(ss))
	mux.Handle("/secure/home", handleSecureHome(ss))
	mux.Handle("/templates", handleTemplates(tc))
	mux.Handle("/bootstrap", handleBootstrapExample())
	mux.Handle("/static/", StaticHandler("/static", "pkg/webapp/example/main/web/static/"))
	log.Fatal(http.ListenAndServe(":8080", mux))

}

func StaticHandler(prefix, path string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
}

func handleIndex(t *webapp.TemplateCache) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.html", map[string]interface{}{})
	}
	return http.HandlerFunc(fn)
}

func handleLogin(t *webapp.TemplateCache, ss *webapp.SessionStore) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			t.ExecuteTemplate(w, "login.html", map[string]interface{}{})
			return
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}
			u, p := r.Form.Get("username"), r.Form.Get("password")
			if u == "admin" && p == "admin" {
				ss.NewSession(w, r)
				http.Redirect(w, r, "/secure/home", http.StatusTemporaryRedirect)
				return
			}
			http.Redirect(w, r, "/login?error=invalid", http.StatusTemporaryRedirect)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func handleSecureHome(ss *webapp.SessionStore) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, ok := ss.CurrentUser(r)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "this is my secure home")
		return
	}
	return http.HandlerFunc(fn)
}

func handleLogout(ss *webapp.SessionStore) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ss.EndSession(w, r)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
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

func handleBootstrapExample() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/webapp/example/main/web/templates/bootstrap-template.html")
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
