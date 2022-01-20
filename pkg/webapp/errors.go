package webapp

import (
	"html/template"
	"net/http"
	"strconv"
)

var defaultErrTmpl = template.Must(template.New("error").Parse(`THIS IS THE DEFAULT ERROR TEMPLATE`))

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
