package webapp

import "net/http"

func StaticHandler(prefix, path string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
}

func AuthHandler(dao DataAccesser) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// do stuff
	}
	return http.HandlerFunc(fn)
}
