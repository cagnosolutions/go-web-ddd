package webappr

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp"
	"net/http"
	"strconv"
)

var httpHandlerRoutes = []struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}{
	{"GET", "/", home},
	{"GET", "/contact", contact},
	{"GET", "/api/widgets", apiGetWidgets},
	{"POST", "/api/widgets", apiCreateWidget},
	{"POST", "/api/widgets/:slug", apiUpdateWidget},
	{"POST", "/api/widgets/:slug/parts", apiCreateWidgetPart},
	{"POST", "/api/widgets/:slug/parts/:id/update", apiUpdateWidgetPart},
	{"POST", "/api/widgets/:slug/parts/:id/delete", apiDeleteWidgetPart},
	{"GET", "/:slug", widget},
	{"GET", "/:slug/admin", widgetAdmin},
	{"POST", "/:slug/image", widgetImage},
}

func Makewebappmuxer() *webapp.Muxer {
	m := webapp.NewMuxer(&webapp.MuxerConfig{
		StaticHandler: nil,
		ErrHandler:    nil,
		MetricsOn:     false,
		Logging:       webapp.LevelOff,
	})
	for _, route := range httpHandlerRoutes {
		m.Handle(route.Method, route.Path, route.Handler)
	}
	return m
}

func getField(r *http.Request, key string) string {
	return r.FormValue(key)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home\n")
}

func contact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "contact\n")
}

func apiGetWidgets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "apiGetWidgets\n")
}

func apiCreateWidget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "apiCreateWidget\n")
}

func apiUpdateWidget(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, "slug")
	fmt.Fprintf(w, "apiUpdateWidget %s\n", slug)
}

func apiCreateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, "slug")
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", slug)
}

func apiUpdateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, "slug")
	id, _ := strconv.Atoi(getField(r, "id"))
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", slug, id)
}

func apiDeleteWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, "slug")
	id, _ := strconv.Atoi(getField(r, "id"))
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", slug, id)
}

func widget(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, "slug")
	fmt.Fprintf(w, "widget %s\n", slug)
}

func widgetAdmin(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, "slug")
	fmt.Fprintf(w, "widgetAdmin %s\n", slug)
}

func widgetImage(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, "slug")
	fmt.Fprintf(w, "widgetImage %s\n", slug)
}
