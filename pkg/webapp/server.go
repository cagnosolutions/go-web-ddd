package webapp

import (
	"net/http"
)

type ServerConfig struct {
	ListenAddr     string
	DefaultHandler http.Handler
}

type ApplicationConfig struct {
	*SessionConfig
	*BasicAuthUser
	*ServerConfig
}

func checkConf(conf *ApplicationConfig) {

}

type Application struct {
	conf *ApplicationConfig
	sess *SessionStore
	auth *BasicAuthUser
	serv *http.Server
}

func NewApplication(conf *ApplicationConfig) *Application {
	checkConf(conf)
	return &Application{
		conf: conf,
		sess: NewSessionStore(conf.SessionConfig),
		auth: NewBasicAuthUser(),
		serv: &http.Server{
			Addr:    conf.ServerConfig.ListenAddr,
			Handler: conf.ServerConfig.DefaultHandler,
		},
	}
}

func (app *Application) ListenAndServe() error {
	return app.serv.ListenAndServe()
}
