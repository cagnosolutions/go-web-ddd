package webapp

import (
	"strings"
)

const pathSeperator = "/"

type Path struct {
	Path string
	ID   string
}

func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, pathSeperator)
	s := strings.Split(p, pathSeperator)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], pathSeperator)
	}
	return &Path{
		Path: p,
		ID:   id,
	}
}

func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
