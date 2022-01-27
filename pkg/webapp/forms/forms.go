package forms

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

var formTemplate = template.Must(template.New("form").Funcs(template.FuncMap{
	"lower": strings.ToLower,
	"title": strings.Title,
	"type":  FieldTypeString,
}).Parse(formStr))

var formStr = `<div class="row row-pad">
                <br>
                <legend>{{ .Name }}</legend>
                <hr>
                <form id="login-form" action="{{ .Action }}" method="post" novalidate="novalidate" autocomplete="off">
				{{ range .Fields }}
		<div class="mb-3">
			<label for="{{ lower .Name }}" class="form-label">{{ title .Name }}</label>
			<input type="{{ type .Type }}" class="form-control" name="{{ lower .Name }}" id="{{ lower .ID }}" aria-describedby="{{ lower .ID }}-help">
    		{{ if ne .HelpText "" }}<div id="{{ lower .ID }}-help" class="form-text">{{ .HelpText }}</div>{{ end }}
		</div>
				{{ end }}
                    <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                        <button type="submit" class="btn btn-success me-md-2">{{ .SubmitText }}</button>
						{{ if .HasCancel }}
							<button type="cancel" class="btn btn-danger me-md-2">Cancel</button>
                    	{{ end }}
					</div>
                </form>
            </div>`

type Field struct {
	ID          string
	Name        string
	Kind        FieldType
	Value       string
	Label       string
	HelpText    string
	MinLen      int
	MaxLen      int
	Placeholder string
	ErrorMsg    string
	Required    bool
	Disabled    bool
}

func MakeFormFromStruct(p interface{}) (string, error) {
	return "", nil
}

func (f Field) String() string {
	strs := []string{
		fmt.Sprintf(`<label for="%s" class="form-label">%s</label>`,
			strings.ToLower(f.Name), strings.ToTitle(f.Name)),
		fmt.Sprintf(`<input type="%s" class="form-control" name="%s" id="%s" aria-describedby="%s-help">`,
			FieldTypeString(f.Kind), strings.ToLower(f.Name), strings.ToLower(f.ID), strings.ToLower(f.ID)),
		fmt.Sprintf(`<div id="%s-help" class="form-text">%s</div>`,
			strings.ToLower(f.ID), f.HelpText),
	}
	return strings.Join(strs, "\n")
}

type FormField interface {
	Type() FieldType
}

type Form struct {
	Name       string
	Action     string
	Fields     []FormField
	SubmitText string
	HasCancel  bool
}

func (f *Form) String() string {
	buf := new(bytes.Buffer)
	err := formTemplate.Execute(buf, f)
	if err != nil {
		panic("form template panic:" + err.Error())
	}
	return buf.String()
}

func MakeForm(name, action, submitText string, hasCancel bool, fields ...FormField) *Form {
	if submitText == "" {
		submitText = "Submit"
	}
	form := &Form{
		Name:       name,
		Action:     action,
		SubmitText: submitText,
		HasCancel:  hasCancel,
		Fields:     make([]FormField, len(fields)),
	}
	for i := range form.Fields {
		form.Fields[i] = fields[i]
	}
	return form
}

type FieldType int

const (
	TypeHidden FieldType = iota
	TypeText
	TypeEmail
	TypeInteger
)

func FieldTypeString(f FieldType) string {
	switch f {
	case TypeHidden:
		return "hidden"
	case TypeText:
		return "text"
	case TypeEmail:
		return "email"
	case TypeInteger:
		return "integer"
	default:
		return "text"
	}
}

type TextField Field

func (f TextField) Type() FieldType {
	return TypeText
}

type EmailField Field

func (f EmailField) Type() FieldType {
	return TypeEmail
}

type IntegerField Field

func (f IntegerField) Type() FieldType {
	return TypeEmail
}

type _BooleanField Field
type _ChoiceField Field
type _DateField Field
type _DateTimeField Field
type _FileField Field
type _ImageField Field
type _NumberField Field
type _DecimalField Field
type _MultipleChoiceField Field
type _ Field
type _ Field
type _ Field
type _ Field
type _ Field
type _ Field

type Text Field
type Email Field
type Password Field
type Tel Field
type Date Field
type Hidden Field
type File Field

type Checkbox Field
type Color Field

type DatetimeLocal Field

type Image Field
type Month Field
type Number Field

type Radio Field
type Range Field

type Search Field

type Reset Field
type Button Field

type Submit Field

type Time Field
type Url Field
type Week Field
