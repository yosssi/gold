package gold

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

const (
	unicodeTab            = 9
	unicodeSpace          = 32
	unicodeDoubleQuote    = 34
	indentTop             = 0
	extendsBlockTokensLen = 2
	goldExtension         = ".gold"
)

// A generator represents an HTML generator.
type Generator struct {
	cache       bool
	templates   map[string]*template.Template
	gtemplates  map[string]*Template
	helperFuncs template.FuncMap
}

// ParseFile parses a Gold template file and returns an HTML template.
func (g *Generator) ParseFile(path string) (*template.Template, error) {
	return g.generateTemplate(path, nil)
}

//Set Helpers
func (g *Generator) SetHelpers(helperFuncs template.FuncMap) {
	g.helperFuncs = helperFuncs
}

// ParseString parses a Gold template string and returns an HTML template.
func (g *Generator) ParseString(stringTemplates map[string]string, name string) (*template.Template, error) {
	return g.generateTemplate(name, stringTemplates)
}

// generateTemplate parses a Gold template and returns an HTML template.
func (g *Generator) generateTemplate(path string, stringTemplates map[string]string) (*template.Template, error) {
	if g.cache {
		if tpl, prs := g.templates[path]; prs {
			return tpl, nil
		}
	}
	gtpl, err := g.parse(path, stringTemplates)
	if err != nil {
		return nil, err
	}
	html, err := gtpl.Html(stringTemplates)
	if err != nil {
		return nil, err
	}
	tpl := template.New(path)
	tpl.Funcs(g.helperFuncs)
	tpl.Parse(html)
	if err != nil {
		return nil, err
	}
	if g.cache {
		g.templates[path] = tpl
	}
	return tpl, nil
}

// parse parses a Gold template file and returns a Gold template.
func (g *Generator) parse(path string, stringTemplates map[string]string) (*Template, error) {
	if g.cache {
		if tpl, prs := g.gtemplates[path]; prs {
			return tpl, nil
		}
	}
	var s string
	if stringTemplates == nil {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		s = string(b)
	} else {
		s = stringTemplates[path]
	}
	lines := strings.Split(formatLf(s), "\n")
	i, l := 0, len(lines)
	tpl := NewTemplate(path, g)
	for i < l {
		line := lines[i]
		i++
		if empty(line) {
			continue
		}
		if topElement(line) {
			switch {
			case isExtends(line):
				tokens := strings.Split(strings.TrimSpace(line), " ")
				if l := len(tokens); l != extendsBlockTokensLen {
					return nil, errors.New(fmt.Sprintf("The line tokens length is invalid. (expected: %d, actual: %d, line no: %d)", extendsBlockTokensLen, l, i))
				}
				superTplPath := tokens[1]
				var superTpl *Template
				var err error
				if stringTemplates == nil {
					superTpl, err = g.parse(tpl.Dir()+superTplPath+goldExtension, nil)
				} else {
					superTpl, err = g.parse(superTplPath, stringTemplates)
				}
				if err != nil {
					return nil, err
				}
				superTpl.Sub = tpl
				tpl.Super = superTpl
			case tpl.Super != nil && isBlock(line):
				tokens := strings.Split(strings.TrimSpace(line), " ")
				if l := len(tokens); l != extendsBlockTokensLen {
					return nil, errors.New(fmt.Sprintf("The lien tokens length is invalid. (expected: %d, actual: %d, line no: %d)", extendsBlockTokensLen, l, i))
				}
				block := &Block{Name: tokens[1], Template: tpl}
				tpl.AddBlock(block.Name, block)
				if err := appendChildren(block, lines, &i, &l, indentTop, false, ""); err != nil {
					return nil, err
				}
			default:
				e, err := NewElement(line, i, indentTop, nil, tpl, nil)
				if err != nil {
					return nil, err
				}
				tpl.AppendElement(e)
				if err := appendChildren(e, lines, &i, &l, indentTop, e.RawContent, e.Type); err != nil {
					return nil, err
				}
			}
		}
	}
	if g.cache {
		g.gtemplates[path] = tpl
	}
	return tpl, nil
}

// NewGenerator generages a generator and returns it.
func NewGenerator(cache bool) *Generator {
	return &Generator{cache: cache, templates: make(map[string]*template.Template), gtemplates: make(map[string]*Template)}
}

// formatLf returns a string whose line feed codes are replaced with LF.
func formatLf(s string) string {
	return strings.Replace(strings.Replace(s, "\r\n", "\n", -1), "\r", "\n", -1)
}

// indent returns the string's indent.
func indent(s string) int {
	i := 0
	space := false
indentLoop:
	for _, b := range s {
		switch b {
		case unicodeTab:
			i++
		case unicodeSpace:
			if space {
				i++
				space = false
			} else {
				space = true
			}
		default:
			break indentLoop
		}
	}
	return i
}

// empty returns if the string is empty or not.
func empty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// topElement returns if the string is the top element.
func topElement(s string) bool {
	return indent(s) == indentTop
}

// appendChildren fetches the lines and appends child elements to the element.
func appendChildren(parent Container, lines []string, i *int, l *int, parentIndent int, parentRawContent bool, parentType string) error {
	for *i < *l {
		line := lines[*i]
		if empty(line) {
			*i++
			continue
		}
		indent := indent(line)
		switch {
		case parentRawContent || parentType == TypeContent:
			switch {
			case indent < parentIndent+1:
				return nil
			default:
				if err := appendChild(parent, &line, &indent, lines, i, l); err != nil {
					return err
				}
			}
		case parentType == TypeBlock:
			switch {
			case indent < parentIndent+1:
				return nil
			default:
				return errors.New(fmt.Sprintf("The indent of the line %d is invalid. Block element can not have child elements.", *i+1))
			}
		default:
			switch {
			case indent < parentIndent+1:
				return nil
			case indent == parentIndent+1:
				if err := appendChild(parent, &line, &indent, lines, i, l); err != nil {
					return err
				}
			case indent > parentIndent+1:
				return errors.New(fmt.Sprintf("The indent of the line %d is invalid.", *i+1))
			}
		}
	}
	return nil
}

// appendChild appends the child element to the parent element.
func appendChild(parent Container, line *string, indent *int, lines []string, i *int, l *int) error {
	var child *Element
	var err error
	switch p := parent.(type) {
	case *Block:
		child, err = NewElement(*line, *i+1, *indent, nil, nil, p)
	case *Element:
		child, err = NewElement(*line, *i+1, *indent, p, nil, nil)
	}
	if err != nil {
		return err
	}
	parent.AppendChild(child)
	*i++
	err = appendChildren(child, lines, i, l, child.Indent, child.RawContent, child.Type)
	if err != nil {
		return err
	}
	return nil
}

// isExtends returns if the line's prefix is "extends" or not.
func isExtends(line string) bool {
	return strings.HasPrefix(line, "extends ") || line == "extends"
}

// isBlock returns if the line's prefix is "block" or not.
func isBlock(line string) bool {
	return strings.HasPrefix(line, "block ") || line == "block"
}
