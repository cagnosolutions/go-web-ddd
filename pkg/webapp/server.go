package webapp

import (
	"net/http"
)

type ApplicationConfig struct {
	*SessionConfig
	Addr    string
	Handler http.Handler
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
			Addr:    conf.Addr,
			Handler: conf.Handler,
		},
	}
}

func (app *Application) ListenAndServe() error {
	return app.serv.ListenAndServe()
}
