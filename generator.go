package gold

import (
	"errors"
	"fmt"
	"github.com/yosssi/gold/template"
	"io/ioutil"
	"strings"
)

const (
	unicodeTab         = 9
	unicodeSpace       = 32
	unicodeDoubleQuote = 34
	indentTop          = 0
)

// A generator represents an HTML generator.
type generator struct {
	cache     bool
	templates map[string]*template.Template
}

// Html parses a template and returns an html string.
func (g *generator) Html(path string, data interface{}) (string, error) {
	tpl, err := g.parse(path)
	if err != nil {
		return "", err
	}
	html, err := tpl.Html()
	if err != nil {
		return "", err
	}
	return html, nil
}

// parse parses a Gold template file and returns a template.
func (g *generator) parse(path string) (*template.Template, error) {
	if g.cache {
		if tpl, prs := g.templates[path]; prs {
			return tpl, nil
		}
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(formatLf(string(b)), "\n")
	i, l := 0, len(lines)
	tpl := template.NewTemplate(path)
	for i < l {
		line := lines[i]
		i++
		if empty(line) {
			continue
		}
		if topElement(line) {
			e, err := template.NewElement(line, i, indentTop, nil, tpl)
			if err != nil {
				return nil, err
			}
			tpl.AppendElement(e)
			err = appendChildren(e, lines, &i, &l)
			if err != nil {
				return nil, err
			}
		}
	}
	if g.cache {
		g.templates[path] = tpl
	}
	return tpl, nil
}

// NewGenerator generages a generator and returns it.
func NewGenerator(cache bool) generator {
	return generator{cache: cache, templates: make(map[string]*template.Template)}
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

// empty returns if the string is empty.
func empty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// topElement returns if the string is the top element.
func topElement(s string) bool {
	return indent(s) == indentTop
}

// appendChildren fetches the lines and appends child elements to the element.
func appendChildren(parent *template.Element, lines []string, i *int, l *int) error {
	for *i < *l {
		line := lines[*i]
		if empty(line) {
			*i++
			continue
		}
		indent := indent(line)
		switch {
		case parent.Tag == "script" || parent.Tag == "style" || parent.Type == template.TypeScriptStyleContent:
			switch {
			case indent < parent.Indent+1:
				return nil
			default:
				if err := appendChild(parent, &line, &indent, lines, i, l); err != nil {
					return err
				}
			}
		default:
			switch {
			case indent < parent.Indent+1:
				return nil
			case indent == parent.Indent+1:
				if err := appendChild(parent, &line, &indent, lines, i, l); err != nil {
					return err
				}
			case indent > parent.Indent+1:
				return errors.New(fmt.Sprintf("The indent of the line %d is invalid.", *i+1))
			}
		}
	}
	return nil
}

// appendChild appends the child element to the parent element.
func appendChild(parent *template.Element, line *string, indent *int, lines []string, i *int, l *int) error {
	child, err := template.NewElement(*line, *i+1, *indent, parent, nil)
	if err != nil {
		return err
	}
	parent.AppendChild(child)
	*i++
	err = appendChildren(child, lines, i, l)
	if err != nil {
		return err
	}
	return nil
}
