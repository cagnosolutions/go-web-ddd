package webapp

import "strings"

const pathSeperator = "/:"

type path struct {
	Path string
	ID   string
}

func NewPath(p string) *path {
	var id string
	p = strings.Trim(p, pathSeperator)
	s := strings.Split(p, pathSeperator)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], pathSeperator)
	}
	return &path{
		Path: p,
		ID:   id,
	}
}

func (p *path) HasID() bool {
	return len(p.ID) > 0
}
