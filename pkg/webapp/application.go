package webapp

import (
	"net/http"
	"time"
)

type WebServerConfig struct {
	ListenAddr     string
	DefaultHandler http.Handler
}

type WebServer struct {
	serv *http.Server
}

func NewWebServer(conf *WebServerConfig) *WebServer {
	return &WebServer{
		serv: &http.Server{
			Addr:    conf.ListenAddr,
			Handler: conf.DefaultHandler,
		},
	}
}

type ApplicationConfig struct {
	*BasicAuthUser
	*SessionConfig
	*TemplateConfig
	*WebServerConfig
}

var defaultApplicationConfig = &ApplicationConfig{
	BasicAuthUser: NewBasicAuthUser(),
	SessionConfig: &SessionConfig{
		SessionID: "go_sess_id",
		Domain:    "localhost",
		Timeout:   time.Duration(15) * time.Minute,
	},
	TemplateConfig: &TemplateConfig{
		BasePattern:   "web/templates/*.html",
		ExtraPatterns: []string{"web/templates/stubs/*.html", "web/templates/misc/*.html"},
		FuncMap:       nil,
	},
	WebServerConfig: &WebServerConfig{
		ListenAddr:     ":8080",
		DefaultHandler: nil,
	},
}

type Application struct {
	conf *ApplicationConfig
	auth *BasicAuthUser
	sess *SessionStore
	tmpl *TemplateCache
	serv *WebServer
}

func checkConf(conf *ApplicationConfig) {

}

func NewApplication(conf *ApplicationConfig) *Application {
	checkConf(conf)
	return &Application{
		conf: conf,
		auth: conf.BasicAuthUser,
		sess: NewSessionStore(conf.SessionConfig),
		tmpl: NewTemplateCache(conf.TemplateConfig),
		serv: NewWebServer(conf.WebServerConfig),
	}
}

func (app *Application) SessionStore() *SessionStore {
	return app.sess
}

func (app *Application) TemplateCache() *TemplateCache {
	return app.tmpl
}

func (app *Application) WebServer() *WebServer {
	return app.serv
}

func (app *Application) SetServerHandler(h http.Handler) {
	app.serv.serv.Handler = h
}

func (app *Application) ListenAndServe() error {
	return app.serv.serv.ListenAndServe()
}
