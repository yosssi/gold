package gold

import (
	"fmt"
	"github.com/yosssi/gold/template"
	"io/ioutil"
	"strings"
)

const (
	unicodeTab         = 9
	unicodeDoubleQuote = 34
	indentTop          = 0
)

// A generator represents an HTML generator.
type generator struct {
	cache     bool
	templates []template.Template
}

// Html parses a template and returns an html string.
func (g *generator) Html(path string, data interface{}) (string, error) {
	tpl, err := g.parse(path)
	if err != nil {
		return "", err
	}
	fmt.Println(tpl)
	return "", nil
}

// parse parses a Gold template file and returns a template.
func (g *generator) parse(path string) (template.Template, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return template.Template{}, err
	}
	lines := strings.Split(formatLf(string(b)), "\n")
	i, l := 0, len(lines)
	tpl := template.Template{}
	for i < l {
		line := lines[i]
		i++
		if empty(line) {
			continue
		}
		if topElement(line) {
		}
	}
	return tpl, nil
}

// NewGenerator generages a generator and returns it.
func NewGenerator(cache bool) generator {
	return generator{cache: cache}
}

// formatLf returns a string whose line feed codes are replaced with LF.
func formatLf(s string) string {
	return strings.Replace(strings.Replace(s, "\r\n", "\n", -1), "\r", "\n", -1)
}

// indent returns the string's indent.
func indent(s string) int {
	i := 0
	for _, b := range s {
		if b == unicodeTab {
			i++
		} else {
			break
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
