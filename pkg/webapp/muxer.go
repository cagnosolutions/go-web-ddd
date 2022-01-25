package webapp

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"path"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type muxEntry struct {
	method  string
	pattern string
	handler http.Handler
}

func (m muxEntry) String() string {
	if m.method == http.MethodGet {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"%s\">%s</a>", m.method, m.pattern, m.pattern)
	}
	if m.method == http.MethodPost {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;%s", m.method, m.pattern)
	}
	if m.method == http.MethodPut {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;&nbsp;%s", m.method, m.pattern)
	}
	if m.method == http.MethodDelete {
		return fmt.Sprintf("[%s]&nbsp;%s", m.method, m.pattern)
	}
	return fmt.Sprintf("[%s]&nbsp;%s", m.method, m.pattern)
}

func (s *Muxer) Len() int {
	return len(s.es)
}

func (s *Muxer) Less(i, j int) bool {
	return s.es[i].pattern < s.es[j].pattern
}

func (s *Muxer) Swap(i, j int) {
	s.es[j], s.es[i] = s.es[i], s.es[j]
}

func (s *Muxer) Search(x string) int {
	return sort.Search(len(s.es), func(i int) bool {
		return s.es[i].pattern >= x
	})
}

type MuxerConfig struct {
	StaticHandler http.Handler
	ErrHandler    http.Handler
	MetricsOn     bool
	Logging       int
}

var defaultMuxerConfig = &MuxerConfig{
	StaticHandler: DefaultMuxerStaticHandler("web/static/"),
	ErrHandler:    DefaultMuxerErrorHandler(),
	MetricsOn:     false,
	Logging:       LevelInfo,
}

func checkMuxerConfig(conf *MuxerConfig) {
	if conf == nil {
		conf = defaultMuxerConfig
	}
}

type Muxer struct {
	conf        *MuxerConfig
	lock        sync.RWMutex
	em          map[string]muxEntry
	es          []muxEntry
	logger      *Logger
	withLogging bool
}

// cleanPath returns the canonical path for p, eliminating . and .. elements
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}

func NewMuxer(conf *MuxerConfig) *Muxer {
	checkMuxerConfig(conf)
	mux := &Muxer{
		conf: conf,
		em:   make(map[string]muxEntry),
		es:   make([]muxEntry, 0),
	}
	if conf.Logging < LevelOff {
		mux.logger = NewLogger(conf.Logging)
		mux.withLogging = true
	}
	if conf.StaticHandler != nil {
		mux.Get("/static/", conf.StaticHandler)
	}
	if conf.ErrHandler != nil {
		mux.Get("/error/", conf.ErrHandler)
	}
	if conf.MetricsOn {
		mux.Get("/metrics", mux.info())
	}
	return mux
}

func (s *Muxer) Handle(method string, pattern string, handler http.Handler) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := s.em[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}
	entry := muxEntry{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
	s.em[pattern] = entry
	if pattern[len(pattern)-1] == '/' {
		s.es = appendSorted(s.es, entry)
	}
}

func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
	n := len(es)
	i := sort.Search(n, func(i int) bool {
		return len(es[i].pattern) < len(e.pattern)
	})
	if i == n {
		return append(es, e)
	}
	// we now know that i points at where we want to insert
	es = append(es, muxEntry{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])      // Move shorter entries down
	es[i] = e
	return es
}

func (s *Muxer) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	s.Handle(method, pattern, http.HandlerFunc(handler))
}

func (s *Muxer) Forward(oldpattern string, newpattern string) {
	s.Handle(http.MethodGet, oldpattern, http.RedirectHandler(newpattern, http.StatusTemporaryRedirect))
}

func (s *Muxer) Get(pattern string, handler http.Handler) {
	s.Handle(http.MethodGet, pattern, handler)
}

func (s *Muxer) Post(pattern string, handler http.Handler) {
	s.Handle(http.MethodPost, pattern, handler)
}

func (s *Muxer) Put(pattern string, handler http.Handler) {
	s.Handle(http.MethodPut, pattern, handler)
}

func (s *Muxer) Delete(pattern string, handler http.Handler) {
	s.Handle(http.MethodDelete, pattern, handler)
}

func (s *Muxer) Static(pattern string, path string) {
	staticHandler := http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
	s.Handle(http.MethodGet, pattern, staticHandler)
}

func (s *Muxer) GetEntries() []string {
	return s.getEntries()
}

func (s *Muxer) getEntries() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	var entries []string
	for _, entry := range s.em {
		entries = append(entries, fmt.Sprintf("%s %s\n", entry.method, entry.pattern))
	}
	return entries
}

// match attempts to locate a handler on a handler map given a
// path string; most-specific (longest) pattern wins
func (s *Muxer) match(path string) (string, string, http.Handler) {
	// first, check for exact match
	e, ok := s.em[path]
	if ok {
		return e.method, e.pattern, e.handler
	}
	// then, check for longest valid match. mux.es
	// contains all patterns that end in "/" sorted
	// from longest to shortest
	for _, e = range s.es {
		if strings.HasPrefix(path, e.pattern) {
			return e.method, e.pattern, e.handler
		}
	}
	return "", "", nil
}

func (s *Muxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, _, h := s.match(r.URL.Path)
	if m != r.Method {
		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
		})
	}
	if h == nil {
		h = http.NotFoundHandler()
	}
	if s.withLogging {
		// if logging is configured, then log, otherwise skip
		h = s.requestLogger(h)
	}
	h.ServeHTTP(w, r)
}

func (s *Muxer) info() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var data []string
		data = append(data, fmt.Sprintf("<h3>Registered Routes (%d)</h3>", len(s.em)))
		for _, entry := range s.em {
			data = append(data, entry.String())
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i] < data[j]
		})
		s.ContentType(w, ".html")
		_, err := fmt.Fprintf(w, strings.Join(data, "<br>"))
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
		return
	}
	return http.HandlerFunc(fn)
}

func (s *Muxer) ContentType(w http.ResponseWriter, content string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	ct := mime.TypeByExtension(content)
	if ct == "" && s.withLogging {
		s.logger.Error("Error, incompatible content type!\n")
		return
	}
	w.Header().Set("Content-Type", ct)
	return
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.data.status = statusCode
}

func (s *Muxer) requestLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				s.logger.Error("err: %v, trace: %s\n", err, debug.Stack())
			}
		}()
		lrw := loggingResponseWriter{
			ResponseWriter: w,
			data: &responseData{
				status: 200,
				size:   0,
			},
		}
		next.ServeHTTP(&lrw, r)
		if 400 <= lrw.data.status && lrw.data.status <= 599 {
			str, args := logStr(lrw.data.status, r)
			s.logger.Error(str, args...)
			return
		}
		str, args := logStr(lrw.data.status, r)
		s.logger.Info(str, args...)
		return
	}
	return http.HandlerFunc(fn)
}

func logStr(code int, r *http.Request) (string, []interface{}) {
	return "# %s - - [%s] \"%s %s %s\" %d %d\n", []interface{}{
		r.RemoteAddr,
		time.Now().Format(time.RFC1123Z),
		r.Method,
		r.URL.EscapedPath(),
		r.Proto,
		code,
		r.ContentLength,
	}
}

func StaticHandler(prefix, path string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
}

func DefaultMuxerStaticHandler(path string) http.Handler {
	return http.StripPrefix("/static", http.FileServer(http.Dir(path)))
}

func DefaultMuxerErrorHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		p := NewPath(r.URL.Path)
		if p.HasID() {
			code, err := strconv.Atoi(p.ID)
			if err != nil {
				code := http.StatusExpectationFailed
				http.Error(w, http.StatusText(code), code)
				return
			}
			err = defaultErrTmpl.Execute(w, struct {
				ErrorCode     int
				ErrorText     string
				ErrorTextLong string
			}{
				ErrorCode:     code,
				ErrorText:     http.StatusText(code),
				ErrorTextLong: HTTPCodesLongFormat[code],
			})
			if err != nil {
				code := http.StatusExpectationFailed
				http.Error(w, http.StatusText(code), code)
				return
			}
		}
	}
	return http.HandlerFunc(fn)
}

var defaultErrTmpl = template.Must(template.New("error.html").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no"/>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <!--[if lt IE 9]>
    <script src="//html5shim.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
    <link href='//fonts.googleapis.com/css?family=Lato:100,300,400,700,900,100italic,300italic,400italic,700italic,900italic'
          rel='stylesheet' type='text/css'>
    <title>Ooops, something went wrong!</title>
</head>

<body>

<!-- navigation -->
<div class="container">
    <nav class="navbar fixed-top navbar-expand-lg navbar-light bg-light">
        <div class="container">
            <a class="navbar-brand" href="#">Ooops, something went wrong!</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href="/back">Take Me Back!</a>
                </div>
            </div>
        </div>
    </nav>
</div>
<div class="navbar-pad"></div>
<!-- navigation -->

<!-- main section -->
<section class="container">
    <div class="fs-1 fw-light text-center">
        <div class="position-absolute top-50 start-50 translate-middle">
            <span>{{ .ErrorCode }}</span>&nbsp;<span class="text-muted fst-italic">{{ .ErrorText }}</span>
            <br>
            <span class="fs-3 fst-italic fw-lighter text-mutex">
                {{ .ErrorTextLong }}
            </span>
            <br>
            <button onclick="history.go(-1)" type="button" class="btn btn-secondary">Please, take me back!</button>
        </div>
    </div>
</section>
<!-- main section -->

<!-- scripts -->
<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.10.2/dist/umd/popper.min.js" integrity="sha384-7+zCNj/IqJ95wo16oMtfsKbZ9ccEh31eOz1HGyDuCQ6wgnyJNSYdrPa03rtR1zdB" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.min.js" integrity="sha384-QJHtvGhmr9XOIpI6YVutG+2QOK9T+ZnN4kzFN1RtK3zEFEIsxhlmWl5/YESvpZ13" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/just-validate@3.3.1/dist/just-validate.production.min.js" crossorigin="anonymous"></script>
<!-- scripts -->

</body>

<!-- footer -->
<nav class="navbar fixed-bottom navbar-light bg-light">
    <div class="container">
        <a class="navbar-brand" href="#"></a>
        <span class="ms-auto">© Some Company 2021-Present</span>
    </div>
</nav>
<!-- footer -->

</html>`))
