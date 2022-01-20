package webapp

import (
	"html/template"
	"net/http"
	"os"
)

type TemplateConfig struct {
	BasePattern   string
	ExtraPatterns []string
	FuncMap       template.FuncMap
}

type TemplateCache struct {
	*TemplateConfig
	FuncMap       template.FuncMap
	t             *template.Template
	basePattern   string
	extraPatterns []string
}

func initTemplates(pattern string, funcMap template.FuncMap) *template.Template {
	return template.Must(template.New("*").Funcs(funcMap).ParseGlob(pattern))
}

func NewTemplateCache(conf *TemplateConfig) *TemplateCache {
	tc := new(TemplateCache)
	if conf.FuncMap == nil {
		conf.FuncMap = template.FuncMap{}
	}
	tc.FuncMap = conf.FuncMap
	tc.basePattern = conf.BasePattern
	tc.t = initTemplates(tc.BasePattern, tc.FuncMap)
	if conf.ExtraPatterns != nil {
		for _, pattern := range conf.ExtraPatterns {
			tc.ParseGlob(pattern)
		}
	}
	return tc
}

func (tc *TemplateCache) ParseGlob(pattern string) {
	t, err := tc.t.Funcs(tc.FuncMap).ParseGlob(pattern)
	if err != nil {
		panic(err)
	}
	tc.t = t
	tc.extraPatterns = append(tc.extraPatterns, pattern)
}

func (tc *TemplateCache) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	err := tc.t.ExecuteTemplate(w, name, data)
	if err != nil {
		//code := http.StatusExpectationFailed
		//http.Error(w, http.StatusText(code), code)
		http.RedirectHandler("/error/417", http.StatusTemporaryRedirect)
		return
	}
}

func (tc *TemplateCache) DefinedTemplates() string {
	return tc.t.DefinedTemplates()
}

func (tc *TemplateCache) Lookup(name string) *template.Template {
	return tc.t.Lookup(name)
}

func (tc *TemplateCache) ReloadTemplates() {
	tc.t = nil
	tc.t = initTemplates(tc.basePattern, tc.FuncMap)
	for i := range tc.extraPatterns {
		t, err := tc.t.Funcs(tc.FuncMap).ParseGlob(tc.extraPatterns[i])
		if err != nil {
			panic(err)
		}
		tc.t = t
	}
}

func FileHasChanged(file string, lastModTime int64) (int64, bool) {
	fi, err := os.Stat(file)
	if err != nil {
		return -1, false
	}
	modTime := fi.ModTime().Unix()
	if modTime > lastModTime {
		return modTime, true
	}
	return modTime, false
}
