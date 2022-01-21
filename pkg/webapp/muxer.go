package webapp

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"
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
	StaticFiles http.Handler
	ErrHandler  http.Handler
	MetricsOn   bool
	Logging     int
}

var defaultMuxerConfig = &MuxerConfig{
	StaticFiles: StaticFileHandler("/static/", "/static", "web/static/"),
	ErrHandler:  ErrorHandler(nil),
	MetricsOn:   false,
	Logging:     LevelInfo,
}

func StaticFileHandler(pattern, prefix, path string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.Handle(pattern,
			http.StripPrefix(prefix,
				http.FileServer(http.Dir(path))))
	}
	return http.HandlerFunc(fn)
}

func checkMuxerConfig(conf *MuxerConfig) {
	if conf == nil {
		conf = defaultMuxerConfig
	}
}

type Muxer struct {
	lock        sync.Mutex
	conf        *MuxerConfig
	logger      *Logger
	withLogging bool
	em          map[string]muxEntry
	es          []muxEntry
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
	if conf.MetricsOn {
		mux.Get("/api/metrics", mux.info())
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
	//s.routes.Put(entry)
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

func (s *Muxer) getEntries() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	var entries []string
	for _, entry := range s.em {
		entries = append(entries, fmt.Sprintf("%s %s\n", entry.method, entry.pattern))
	}
	return entries
}

func (s *Muxer) match(path string) (string, string, http.Handler) {
	e, ok := s.em[path]
	if ok {
		return e.method, e.pattern, e.handler
	}
	for _, e = range s.es {
		if strings.HasPrefix(path, e.pattern) {
			return e.method, e.pattern, e.handler
		}
	}
	return "", "", nil
}

func (s *Muxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, _, h := s.match(r.URL.Path)
	if m != r.Method || h == nil {
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

func ErrorHandler(errTmpl *template.Template) http.Handler {
	if errTmpl == nil {
		errTmpl = defaultErrTmpl
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		p := NewPath(r.URL.Path)
		if p.HasID() {
			code, err := strconv.Atoi(p.ID)
			if err != nil {
				code := http.StatusExpectationFailed
				http.Error(w, http.StatusText(code), code)
				return
			}
			err = errTmpl.Execute(w, struct {
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
