package webapp

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

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
