package webapp

import (
	"fmt"
	"html/template"
	"net/http"
)

type TemplateCache struct {
	t       *template.Template
	FuncMap template.FuncMap
}

func NewTemplateCache(pattern string, funcMap template.FuncMap) *TemplateCache {
	tc := new(TemplateCache)
	if funcMap == nil {
		funcMap = template.FuncMap{}
	}
	tc.FuncMap = funcMap
	t, err := template.New("*").Funcs(tc.FuncMap).ParseGlob(pattern)
	if err != nil {
		panic(err)
	}
	tc.t = t
	return tc
}

func (tc *TemplateCache) ParseGlob(pattern string) {
	t, err := tc.t.Funcs(tc.FuncMap).ParseGlob(pattern)
	if err != nil {
		panic(err)
	}
	tc.t = t
}

func (tc *TemplateCache) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	err := tc.t.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		code := http.StatusExpectationFailed
		http.Error(w, http.StatusText(code), code)
	}
}

func (tc *TemplateCache) DefinedTemplates() string {
	return tc.t.DefinedTemplates()
}

func HandleWithTemplate(h http.Handler, t *template.Template) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

	}
	return http.HandlerFunc(fn)
}
