package webapp

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const defaultSessionInMinutes = 300

func Base64Encode(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

func Base64Decode(s string) string {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(fmt.Sprintf("cookie: base64 decoding failed %q", err))
	}
	return string(b)
}

func URLEncode(s string) string {
	return url.QueryEscape(s)
}

func URLDecode(s string) string {
	us, err := url.QueryUnescape(s)
	if err != nil {
		panic(fmt.Sprintf("cookie: query unescape failed %q", err))
	}
	return us
}

func NewCookie(name, value string, expires time.Time, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:       URLEncode(name),
		Value:      Base64Encode(value),
		Path:       "/",
		Domain:     "",
		Expires:    expires,
		RawExpires: "",
		MaxAge:     maxAge,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
}

func GetCookie(r *http.Request, name string) *http.Cookie {
	c, err := r.Cookie(URLEncode(name))
	if err != nil || err == http.ErrNoCookie {
		return nil
	}
	c.Value = Base64Decode(c.Value)
	return c
}

func SetCookie(w http.ResponseWriter, name, value string) {
	c := NewCookie(name, value, time.Now().Add(defaultSessionInMinutes*time.Minute), 0)
	http.SetCookie(w, c)
}

func DelCookie(w http.ResponseWriter, name string) {
	c := NewCookie(name, "", time.Now(), 0)
	c.MaxAge = -1
	http.SetCookie(w, c)
}
