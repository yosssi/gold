package template

import (
	"bytes"
	"strings"
)

// A template represents a Gold template.
type Template struct {
	Path     string
	Elements []*Element
	Super    *Template
	Blocks   map[string]*Block
}

// AppendElement appends the element to the template's elements.
func (t *Template) AppendElement(e *Element) {
	t.Elements = append(t.Elements, e)
}

// Html generates an html and returns it.
func (t *Template) Html() (string, error) {
	var bf bytes.Buffer
	for _, e := range t.Elements {
		err := e.html(&bf)
		if err != nil {
			return "", err
		}
	}
	return bf.String(), nil
}

// Dir returns the template file's directory.
func (t *Template) Dir() string {
	tokens := strings.Split(t.Path, "/")
	l := len(tokens)
	switch {
	case l < 2:
		return "./"
	default:
		return strings.Join(tokens[:l-1], "/") + "/"
	}
}

// NewTemplate generates a new template and returns it.
func NewTemplate(path string) *Template {
	return &Template{Path: path, Blocks: make(map[string]*Block)}
}
