package webapp

import "net/http"

func StaticHandler(prefix, path string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
}
