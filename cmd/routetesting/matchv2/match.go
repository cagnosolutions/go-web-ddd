package matchv2

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	p := r.URL.Path
	if _, ok := match(p, "/"); ok {
		h = get(home)
		h.ServeHTTP(w, r)
		return
	}
	if _, ok := match(p, "/contact"); ok {
		h = get(contact)
		h.ServeHTTP(w, r)
		return
	}
	if _, ok := match(p, "/api/widgets"); ok && r.Method == "GET" {
		h = get(apiGetWidgets)
		h.ServeHTTP(w, r)
		return
	}
	if _, ok := match(p, "/api/widgets"); ok {
		h = post(apiCreateWidget)
		h.ServeHTTP(w, r)
		return
	}
	if v, ok := match(p, "/api/widgets/:slug"); ok {
		h = post(apiWidget{
			v["slug"],
		}.update)
		h.ServeHTTP(w, r)
		return
	}
	if v, ok := match(p, "/api/widgets/:slug/parts"); ok {
		h = post(apiWidget{
			v["slug"],
		}.createPart)
		h.ServeHTTP(w, r)
		return
	}
	if v, ok := match(p, "/api/widgets/:/parts/:/update"); ok {
		id, _ := strconv.Atoi(v["id"])
		h = post(apiWidgetPart{
			v["slug"],
			id,
		}.update)
		h.ServeHTTP(w, r)
		return
	}
	if v, ok := match(p, "/api/widgets/:slug/parts/:id/delete"); ok {
		id, _ := strconv.Atoi(v["id"])
		h = post(apiWidgetPart{
			v["slug"],
			id,
		}.delete)
		h.ServeHTTP(w, r)
		return
	}
	if v, ok := match(p, "/:slug"); ok {
		h = get(widget{
			v["slug"],
		}.widget)
		h.ServeHTTP(w, r)
		return
	}
	if v, ok := match(p, "/:slug/admin"); ok {
		h = get(widget{
			v["slug"],
		}.admin)
		h.ServeHTTP(w, r)
		return
	}
	if v, ok := match(p, "/:slug/image"); ok {
		h = post(widget{
			v["slug"],
		}.image)
		h.ServeHTTP(w, r)
		return
	}
	//fmt.Printf("HANDLING PATH=%q\n", p)
	http.NotFound(w, r)
	return
}

func parse3(path, pattern string) (map[string]string, bool) {
	data := make(map[string]string, 0)
	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case ':':
			// matches till next slash in path
			slash := strings.IndexByte(path, '/')
			if slash < 0 {
				slash = len(path)
			}
			segment := path[:slash]
			path = path[slash:]
			data[segment] = segment
		case path[0]:
			// non-'+' pattern byte must match path byte
			path = path[1:]
		default:
			return nil, false
		}
	}
	return data, path == "" && pattern == ""
}

func parse2(path, pattern string) (map[string]string, bool) {
	//fmt.Printf("1) path=%q (%d), pattern=%q (%d)\n", path, len(path), pattern, len(pattern))
	// first, check for a direct match
	if path == pattern {
		return nil, true
	}
	// next, do a quick suffix check
	if len(path) > 1 && len(pattern) > 1 {
		if path[len(path)-1] == '/' && pattern[len(pattern)-1] != '/' {
			return nil, false
		}
	}
	// next, do some quick omission checks
	//if len(path) > 1 && len(pattern) > 1 {
	//	// check prefix, if they are different
	//	// no possible match can be found here
	//	if path[:2] != pattern[:2] {
	//		return nil, false
	//	}
	//	// if prefix matches, check suffix and if
	//	// there is not a match with the suffix
	//	// then no possible match can be found
	//	if path[len(path)-2] != pattern[len(pattern)-2] {
	//		return nil, false
	//	}
	//}
	// if there are no slugs in the pattern...
	if strings.IndexByte(pattern, ':') == -1 {
		// ...and the length of the pattern and the path are not
		// the same, then there is no way they can be a match
		if len(pattern) != len(path) {
			return nil, false
		}
	}
	// otherwise, there could be a slug and therefore a
	// potential match, but we can easily rule one out
	// if we count the number of separators and if there
	// is a difference between the path and pattern separator
	// count, there is also no possible match.
	if strings.Count(path, "/") != strings.Count(pattern, "/") {
		return nil, false
	}
	// if we get here, that means that we DO NOT have
	// a direct match yet, and it also means that we
	// DO have at least one slug to isolate.
	data := make(map[string]string)
	var hasSlug, hasMatch bool
	for {
		// check for slug in pattern
		if pattern[0] == ':' {
			hasSlug = true
		}
		// get next separator offset
		i := strings.IndexByte(path, '/')
		j := strings.IndexByte(pattern, '/')
		// check if we're at the end
		if i == -1 && j == -1 {
			hasMatch = path == pattern
			//fmt.Printf("i=%d, j=%d, len(path)=%d, len(pattern)=%d, path=%q, pattern=%q, hasSlug=%v, hasMatch=%v\n",
			//	i, j, len(path), len(pattern), path, pattern, hasSlug, hasMatch)
			break
		}
		// otherwise, compare paths (if we have found no slug)
		if path[:i] == pattern[:j] && !hasSlug {
			hasMatch = false
		}
		// check for slug, and add to map
		if hasSlug {
			data[pattern[1:j]] = path[:i]
			hasSlug = false
		}
		// update the path and pattern slices
		path = path[i+1:]
		pattern = pattern[j+1:]
	}
	// return slug map, and match bool
	return data, hasMatch
}

// parse1 reports whether path matches the given pattern
func parse1(path, pattern string) (map[string]string, bool) {
	if path == pattern {
		return nil, path == pattern
	}
	//if strings.IndexByte(pattern, ':') == -1 {
	//	fmt.Printf(">>>>> %q, %q (%v)\n", path, pattern, path == pattern)
	//	return nil, path[1:] == pattern[1:]
	//}
	data := make(map[string]string)
	fn := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != ':'
	}
	ppath := strings.FieldsFunc(path, fn)
	ppatt := strings.FieldsFunc(pattern, fn)
	for i := range ppatt {
		//fmt.Printf(">>>>> %q, %q (%v)\n", ppath[i], ppatt[i], ppath[i] == ppatt[i])
		if ppath[i] != ppatt[i] {
			if ppatt[i][0] == ':' {
				data[ppatt[i][1:]] = ppath[i]
			}
		}
	}
	return data, len(data) > 0
}

// match reports whether path matches the given pattern, which is a
// path with '+' wildcards wherever you want to use a parameter. Path
// parameters are assigned to the pointers in vars (len(vars) must be
// the number of wildcards), which must be of type *string or *int.
func match(path, pattern string) (map[string]string, bool) {
	return parse3(path, pattern)
}

func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Allow", method)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "GET")
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "POST")
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

type apiWidget struct {
	slug string
}

func (h apiWidget) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidget %s\n", h.slug)
}

func (h apiWidget) createPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", h.slug)
}

type apiWidgetPart struct {
	slug string
	id   int
}

func (h apiWidgetPart) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", h.slug, h.id)
}

func (h apiWidgetPart) delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", h.slug, h.id)
}

type widget struct {
	slug string
}

func (h widget) widget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widget %s\n", h.slug)
}

func (h widget) admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetAdmin %s\n", h.slug)
}

func (h widget) image(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetImage %s\n", h.slug)
}
