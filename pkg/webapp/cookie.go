package webapp

import (
	"net/http"
	"time"
)

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

func SaveCookie(w http.ResponseWriter, c *http.Cookie) {
	http.SetCookie(w, c)
}

func DelCookie(w http.ResponseWriter, name string) {
	c := NewCookie(name, "", time.Now(), -1)
	http.SetCookie(w, c)
}
