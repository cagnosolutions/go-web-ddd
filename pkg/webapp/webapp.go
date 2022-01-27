package webapp

import "net/http"

type WebAppConfig struct {
	Templates          *TemplateConfig
	Sessions           *SessionConfig
	Muxer              *MuxerConfig
	Server             *ServerConfig
	AppName            string
	GracefulShutdownOn bool
}

type WebApp struct {
	AuthUser
	*WebAppConfig
	*TemplateCache
	*SessionStore
	*Muxer
	*Server
	onSuccess http.Handler
	onFailure http.Handler
}

func NewWebApp(conf *WebAppConfig) *WebApp {
	if conf == nil {
		panic("WebApp requires configuration")
	}
	app := &WebApp{
		WebAppConfig:  conf,
		TemplateCache: nil,
		SessionStore:  nil,
		Muxer:         nil,
		Server:        nil,
	}
	if conf.Templates != nil {
		app.TemplateCache = NewTemplateCache(conf.Templates)
	}
	if conf.Sessions != nil {
		app.SessionStore = NewSessionStore(conf.Sessions)
		app.AuthUser = NewSystemSessionUser()
	}
	if conf.Muxer != nil {
		app.Muxer = NewMuxer(conf.Muxer)
	}
	if conf.Server != nil {
		app.Server = NewServer(conf.Server)
	}
	if conf.AppName == "" {
		conf.AppName = "Go WebApp"
	}
	if conf.GracefulShutdownOn {
		HandleSignalInterrupt("%q started, and running...\n", conf.AppName)
	}
	return app
}

func (app *WebApp) Redirect(url string) http.Handler {
	return http.RedirectHandler(url, http.StatusTemporaryRedirect)
}

func (app *WebApp) handleLogin() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// handle GET /login
		if r.Method == http.MethodGet {
			app.TemplateCache.ExecuteTemplate(w, "login.html", map[string]interface{}{})
			return
		}
		// handle POST /login
		if r.Method == http.MethodPost {
			// get posted form values
			un := r.FormValue("username")
			pw := r.FormValue("password")
			// attempt to authenticate
			user, ok := app.AuthUser.Authenticate(un, pw)
			if !ok {
				// if authentication failed call onFailure
				app.onFailure.ServeHTTP(w, r)
				return
			}
			// otherwise, start a new session
			sess := app.SessionStore.New()
			sess.Set("role", user.Role)
			sess.Set("username", user.Username)
			app.SessionStore.Save(w, r, sess)
			// call our onSuccess
			app.onSuccess.ServeHTTP(w, r)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func (app *WebApp) HandleLogin() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// reject non-post login calls
		if r.Method != http.MethodPost {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		// get posted form values
		un := r.FormValue("username")
		pw := r.FormValue("password")
		// attempt to authenticate
		user, ok := app.AuthUser.Authenticate(un, pw)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		// otherwise, start a new session
		sess := app.SessionStore.New()
		sess.Set("role", user.Role)
		sess.Set("username", user.Username)
		app.SessionStore.Save(w, r, sess)
		// call onSuccess
		app.onSuccess.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}

func (app *WebApp) HandleRegister() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// reject non-post login calls
		if r.Method != http.MethodPost {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		// get posted form values
		un := r.FormValue("username")
		pw := r.FormValue("password")
		// attempt to authenticate
		user, ok := app.AuthUser.Authenticate(un, pw)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		// otherwise, start a new session
		sess := app.SessionStore.New()
		sess.Set("role", user.Role)
		sess.Set("username", user.Username)
		app.SessionStore.Save(w, r, sess)
		// call onSuccess
		app.onSuccess.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}
