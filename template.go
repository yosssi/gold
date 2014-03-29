package gold

import (
	"bytes"
	"strings"
)

// A template represents a Gold template.
type Template struct {
	Path      string
	Generator *Generator
	Elements  []*Element
	Super     *Template
	Sub       *Template
	Blocks    map[string]*Block
}

// AppendElement appends the element to the template's elements.
func (t *Template) AppendElement(e *Element) {
	t.Elements = append(t.Elements, e)
}

// Html generates an html and returns it.
func (t *Template) Html(stringTemplates map[string]string) (string, error) {
	if t.Super != nil {
		return t.Super.Html(stringTemplates)
	} else {
		var bf bytes.Buffer
		for _, e := range t.Elements {
			err := e.Html(&bf, stringTemplates)
			if err != nil {
				return "", err
			}
		}
		return bf.String(), nil
	}
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

// AddBlock appends the block to the template.
func (t *Template) AddBlock(name string, block *Block) {
	t.Blocks[name] = block
}

// NewTemplate generates a new template and returns it.
func NewTemplate(path string, generator *Generator) *Template {
	return &Template{Path: path, Generator: generator, Blocks: make(map[string]*Block)}
}
