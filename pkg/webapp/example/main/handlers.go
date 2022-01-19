package main

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"net/http"
)

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
			fmt.Println(">>> [GET] >>> LOGIN")
			t.ExecuteTemplate(w, "login.html", map[string]interface{}{})
			return
		case http.MethodPost:
			fmt.Println(">>> [POST] >>> LOGIN")
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
