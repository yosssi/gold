package template

import (
	"bytes"
)

// A template represents a Gold template.
type Template struct {
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
