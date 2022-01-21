package webapp

import "net/http"

type ServerConfig struct {
}

func checkServerConfig(conf *ServerConfig) {

}

type Server struct {
	*ServerConfig
	*http.Server
}

func NewServer(conf *ServerConfig) *Server {
	return &Server{
		ServerConfig: conf,
	}
}
