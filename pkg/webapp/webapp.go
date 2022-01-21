package webapp

type WebAppConfig struct {
	Templates          *TemplateConfig
	Sessions           *SessionConfig
	Muxer              *MuxerConfig
	Server             *ServerConfig
	AppName            string
	GracefulShutdownOn bool
}

func checkWebAppConfig(conf *WebAppConfig) {

}

type WebApp struct {
	*WebAppConfig
	Templates *TemplateCache
	Sessions  *SessionStore
	Muxer     *Muxer
	Server    *Server
}

func NewWebApp(conf *WebAppConfig) *WebApp {
	checkWebAppConfig(conf)
	return &WebApp{
		WebAppConfig: conf,
	}
}
