package webapp

import "net/http"

type Server struct {
	*http.Server
	SessionStore
	AuthUser
}
