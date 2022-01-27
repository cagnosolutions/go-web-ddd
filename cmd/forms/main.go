package main

import (
	"github.com/cagnosolutions/go-web-ddd/pkg/webapp/forms"
)

func main() {
	contactForm, err := forms.MakeFormFromStruct(&ContactUs{})
	if err != nil {
		panic(err)
	}
	// use contact form, it's just a string of HTML
	_ = contactForm
}

// ContactUs is just a normal struct
// with html field tags, required is assumed
// you have to use `required=false` if you
// want to shut it off. html input names and
// id's are also automatically inferred.
type ContactUs struct {
	Name    string `html:"placeholder='your full name'"`
	Subject string `html:"help='100 characters, max'"`
	Message string `html:""`
	Sender  string `html:"help='A valid email address, please'"`
}

func ContactUsForm() *forms.Form {
	name := forms.TextField{
		ID: "full-name", Name: "full-name",
		Placeholder: "Your full name",
		Required:    true,
	}
	subject := forms.TextField{
		ID: "subject", Name: "subject",
		MaxLen:   100,
		HelpText: "100 characters max",
		Required: true,
	}
	message := forms.TextField{
		ID: "subject", Name: "message",
		Required: true,
	}
	sender := forms.EmailField{
		ID: "sender", Name: "sender",
		HelpText: "A valid email address, please",
		Required: true,
	}
	captcha := forms.IntegerField{
		ID: "captcha", Name: "captcha",
		Label:    "2+2=?",
		Required: true,
	}
	return forms.MakeForm("Contact Us", "/contact-us", "Contact Now", false,
		name, subject, message, sender, captcha)
}
