package main

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/example/user"
	"log"
	"net/http"
	"time"
)

var (
	tc *webapp.TemplateCache
	ss *webapp.SessionStore
	ba *webapp.BasicAuthUser
)

func init() {
	// init templates
	tc = webapp.NewTemplateCache(&webapp.TemplateConfig{
		BasePattern:   "pkg/webapp/example/main/web/templates/*.html",
		ExtraPatterns: []string{"pkg/webapp/example/main/web/templates/stubs/*.html"},
		FuncMap:       nil,
	})

	// init session store
	ss = webapp.NewSessionStore(&webapp.SessionConfig{
		SessionID: "sess-id",
		Timeout:   time.Duration(30) * time.Second,
	})

	// init basic auth user
	ba = webapp.NewBasicAuthUser()
	ba.Register("jdoe@example.com", "awesome007", "user")
}

func main() {

	// server
	mux := http.NewServeMux()
	mux.Handle("/error/", webapp.ErrorHandler(tc.Lookup("error.html")))
	mux.Handle("/index", handleIndex(tc))
	mux.Handle("/login", handleLogin(tc, ss, ba))
	mux.Handle("/logout", handleLogout(ss))
	mux.Handle("/sessions", handleSessions(ss))
	mux.Handle("/secure/home", handleSecureHome(ss))
	mux.Handle("/templates", handleTemplates(tc))
	mux.Handle("/bootstrap", handleBootstrapExample())
	mux.Handle("/static/", webapp.StaticHandler("/static", "pkg/webapp/example/main/web/static/"))
	log.Fatal(http.ListenAndServe(":8080", mux))

}

func handleSessions(ss *webapp.SessionStore) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query().Get("sess")
		if v == "new" {
			s := ss.New()
			s.Set("counter", 0)
			ss.Save(w, r, s)
			fmt.Fprintf(w, "got new session (and saved): %s\n", s.ID())
			return
		}
		if v == "get" {
			s, ok := ss.Get(r)
			if !ok {
				fmt.Fprintf(w, "tried to get session but didnt find any")
				return
			}
			c, _ := s.Get("counter")
			s.Set("counter", c.(int)+1)
			ss.Save(w, r, s)
			fmt.Fprintf(w, "got session, count=%d (%d seconds until expire) \n", c, s.ExpiresIn())
			return
		}
		if v == "del" {
			ss.Save(w, r, nil)
			fmt.Fprintf(w, "removed session\n")
			return
		}
		fmt.Fprintf(w, "sessions: %s\n", ss.String())
		return
	}
	return http.HandlerFunc(fn)
}

func loadUser(dao webapp.DataAccesser) (string, http.Handler) {
	// init and wire user
	userController := user.WireUser(dao)

	// add to your main router wherever that is
	return "/user", userController.HandleBase()
}
