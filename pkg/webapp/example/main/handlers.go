package main

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/example/user"
	"log"
	"net/http"
)

func handleIndex(t *webapp.TemplateCache) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.html", map[string]interface{}{})
	}
	return http.HandlerFunc(fn)
}

func handleLogin(t *webapp.TemplateCache, ss *webapp.CookieStore, us *user.UserService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(">>> form: %+v\n", r.Form)
		switch r.Method {
		case http.MethodGet:
			handleLoginGet(t).ServeHTTP(w, r)
		case http.MethodPost:
			handleLoginPost(ss, us).ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func handleLoginGet(t *webapp.TemplateCache) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(">>> [GET] >>> LOGIN")
		t.ExecuteTemplate(w, "login.html", map[string]interface{}{})
	}
	return http.HandlerFunc(fn)
}

func handleLoginPost(ss *webapp.CookieStore, us *user.UserService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		user := r.Form.Get("username")
		pass := r.Form.Get("password")
		u := us.GetUser(user, pass)
		if u == nil {
			fmt.Println("no user found")
		}
		/*
			if u == "admin" && p == "admin" {
				ss.NewSession(w, r)
				http.Redirect(w, r, "/secure/home", http.StatusTemporaryRedirect)
				return
			}
			http.Redirect(w, r, "/login?error=invalid", http.StatusTemporaryRedirect)
			return
		*/
	}
	return http.HandlerFunc(fn)
}

func handleSecureHome(ss *webapp.CookieStore) http.Handler {
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

func handleLogout(ss *webapp.CookieStore) http.Handler {
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
