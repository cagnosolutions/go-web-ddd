package webapp

type WebAppConfig struct {
	Templates          *TemplateConfig
	Sessions           *SessionConfig
	Muxer              *MuxerConfig
	Server             *ServerConfig
	AppName            string
	GracefulShutdownOn bool
}

type WebApp struct {
	*WebAppConfig
	*TemplateCache
	*SessionStore
	*Muxer
	*Server
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
