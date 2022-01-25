package main

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"log"
	"net/http"
)

func handleIndex(t *webapp.TemplateCache) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.html", map[string]interface{}{})
	}
	return http.HandlerFunc(fn)
}

func handleLogin(t *webapp.TemplateCache, ss *webapp.SessionStore, ba *webapp.SystemSessionUser) http.Handler {
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
			user := r.Form.Get("username")
			pass := r.Form.Get("password")
			su, authd := ba.Authenticate(user, pass)
			if !authd {
				t.ExecuteTemplate(w, "login.html", map[string]interface{}{})
				return
			}
			sess := ss.New()
			sess.Set("user", su)
			ss.Save(w, r, sess)
			http.Redirect(w, r, "/secure/home", http.StatusTemporaryRedirect)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func handleSecureHome(ss *webapp.SessionStore) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		sess, ok := ss.Get(r)
		if !ok {
			http.Redirect(w, r, "/error/401", http.StatusTemporaryRedirect)
			return
		}
		usr, _ := sess.Get("user")
		ss.Save(w, r, sess)
		fmt.Fprintf(w, "this is my secure home (session.id=%s, role=%s)\n", sess.ID(), usr)
		return
	}
	return http.HandlerFunc(fn)
}

func handleLogout(ss *webapp.SessionStore) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ss.Save(w, r, nil)
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
