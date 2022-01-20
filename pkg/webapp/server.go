package webapp

import (
	"net/http"
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
	*SessionConfig
	*TemplateConfig
	*BasicAuthUser
	*WebServerConfig
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
		auth: NewBasicAuthUser(),
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
