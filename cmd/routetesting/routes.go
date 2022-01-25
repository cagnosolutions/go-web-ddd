package main

import (
	"net/url"
	"regexp"
)

func MatchStringUsingStrings(str, match string) bool {
	return str == match
}

func initRegex(s []string) []*regexp.Regexp {
	var reg []*regexp.Regexp
	for i := range s {
		r := regexp.MustCompile(s[i])
		reg = append(reg, r)
	}
	return reg
}

func MatchStringUsingRegex(reg *regexp.Regexp, match string) bool {
	return reg.MatchString(match)
}

func Parse(path, s string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(s) {
		switch {
		case j >= len(path):
			if path != "/" && len(path) > 0 && path[len(path)-1] == '/' {
				return p, true
			}
			return nil, false
		case path[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = matcher(path, isBoth, j+1)
			val, _, i = matcher(s, byteParse(nextc), i)
			p.Add(":"+name, val)
		case s[i] == path[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(path) {
		return nil, false
	}
	return p, true
}

// match path with registered handler
func matcher(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

// determine type of byte
func byteParse(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

// test for alpha byte
func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// test for numerical byte
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// test for alpha or numerical byte
func isBoth(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
