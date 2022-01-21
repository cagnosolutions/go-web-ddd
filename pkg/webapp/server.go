package webapp

import "net/http"

type ServerConfig struct {
}

type Server struct {
	*http.Server
}
