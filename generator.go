package gold

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yosssi/gohtml"
)

const (
	unicodeTab            = 9
	unicodeSpace          = 32
	unicodeDoubleQuote    = 34
	indentTop             = 0
	extendsBlockTokensLen = 2
	defaultDelimLeft      = "{{"
	defaultDelimRight     = "}}"
)

// Generator represents an HTML generator.
type Generator struct {
	cache        bool
	templates    map[string]*template.Template
	htmls        map[string]string
	gtemplates   map[string]*Template
	helperFuncs  template.FuncMap
	baseDir      string
	prettyPrint  bool
	debugWriter  io.Writer
	asset        func(string) ([]byte, error)
	assetBaseDir string
	delimLeft    string
	delimRight   string
}

// ParseFile parses a Gold template file and returns an HTML template.
func (g *Generator) ParseFile(path string) (*template.Template, error) {
	tpl, _, err := g.generateTemplate(path, nil, true)
	return tpl, err
}

// ParseFileWithHTML parses a Gold template file and returns an HTML template and HTML source codes.
func (g *Generator) ParseFileWithHTML(path string) (*template.Template, string, error) {
	return g.generateTemplate(path, nil, true)
}

// SetHelpers sets the helperFuncs to the generator.
func (g *Generator) SetHelpers(helperFuncs template.FuncMap) *Generator {
	g.helperFuncs = helperFuncs
	return g
}

// SetBaseDir sets the base directory to the generator.
func (g *Generator) SetBaseDir(baseDir string) *Generator {
	g.baseDir = baseDir
	g.assetBaseDir = baseDir
	return g
}

// SetPrettyPrint sets the prettyPrint to the generator.
func (g *Generator) SetPrettyPrint(prettyPrint bool) *Generator {
	g.prettyPrint = prettyPrint
	return g
}

// SetDebugWriter sets a debugWriter to the generator.
func (g *Generator) SetDebugWriter(debugWriter io.Writer) *Generator {
	g.debugWriter = debugWriter
	return g
}

// SetAsset sets an asset to the generator.
func (g *Generator) SetAsset(asset func(string) ([]byte, error)) *Generator {
	g.asset = asset
	return g
}

// Delims sets the action delimiters to the specified strings
func (g *Generator) Delims(left, right string) *Generator {
	g.delimLeft = left
	g.delimRight = right
	return g
}

// ParseString parses a Gold template string and returns an HTML template.
func (g *Generator) ParseString(stringTemplates map[string]string, name string) (*template.Template, error) {
	tpl, _, err := g.generateTemplate(name, stringTemplates, false)
	return tpl, err
}

// ParseStringWithHTML parses a Gold template string and returns an HTML template and HTML source codes.
func (g *Generator) ParseStringWithHTML(stringTemplates map[string]string, name string) (*template.Template, string, error) {
	return g.generateTemplate(name, stringTemplates, false)
}

// generateTemplate parses a Gold template and returns an HTML template.
func (g *Generator) generateTemplate(path string, stringTemplates map[string]string, addBaseDir bool) (*template.Template, string, error) {
	if g.cache {
		if tpl, prs := g.templates[path]; prs {
			html := g.htmls[path]
			return tpl, html, nil
		}
	}
	gtpl, err := g.parse(path, stringTemplates, addBaseDir)
	if err != nil {
		return nil, "", err
	}
	html, err := gtpl.Html(stringTemplates, nil)
	if err != nil {
		return nil, "", err
	}
	if g.prettyPrint {
		html = gohtml.Format(html)
	}
	if g.debugWriter != nil {
		debugStr := gohtml.AddLineNo(html)
		g.debugWriter.Write([]byte(debugStr + "\n"))
	}
	tpl := template.New(path)
	tpl.Funcs(g.helperFuncs)
	if g.delimLeft != defaultDelimLeft || g.delimRight != defaultDelimRight {
		tpl.Delims(g.delimLeft, g.delimRight)
	}
	_, err = tpl.Parse(html)
	if err != nil {
		return nil, html, err
	}
	if g.cache {
		g.templates[path] = tpl
		g.htmls[path] = html
	}
	return tpl, html, nil
}

// parse parses a Gold template file and returns a Gold template.
func (g *Generator) parse(path string, stringTemplates map[string]string, addBaseDir bool) (*Template, error) {
	if addBaseDir {
		path = Path(g.baseDir, path)
	}
	if g.cache {
		if tpl, prs := g.gtemplates[path]; prs {
			return tpl, nil
		}
	}
	var s string
	if stringTemplates == nil {
		var b []byte
		var err error
		if g.asset == nil {
			b, err = ioutil.ReadFile(path)
		} else {
			b, err = g.asset(assetPath(path, g.baseDir, g.assetBaseDir))
		}
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
					return nil, fmt.Errorf("the line tokens length is invalid. (expected: %d, actual: %d, line no: %d, template: %s, line: %s)", extendsBlockTokensLen, l, i, tpl.Path, strings.TrimSpace(line))
				}
				superTplPath := tokens[1]
				var superTpl *Template
				var err error
				if stringTemplates == nil {
					addBaseDir := true
					if g.baseDir != "" && CurrentDirectoryBasedPath(superTplPath) {
						superTplPath = tpl.Dir() + superTplPath
						addBaseDir = false
					}
					superTpl, err = g.parse(superTplPath+Extension, nil, addBaseDir)
				} else {
					superTpl, err = g.parse(superTplPath, stringTemplates, false)
				}
				if err != nil {
					return nil, err
				}
				superTpl.Sub = tpl
				tpl.Super = superTpl
			case tpl.Super != nil && isBlock(line):
				tokens := strings.Split(strings.TrimSpace(line), " ")
				if l := len(tokens); l != extendsBlockTokensLen {
					return nil, fmt.Errorf("the line tokens length is invalid. (expected: %d, actual: %d, line no: %d, template: %s, line: %s)", extendsBlockTokensLen, l, i, tpl.Path, strings.TrimSpace(line))
				}
				block := &Block{Name: tokens[1], Template: tpl}
				tpl.AddBlock(block.Name, block)
				if err := appendChildren(block, lines, &i, &l, indentTop, false, "", tpl); err != nil {
					return nil, err
				}
			default:
				e, err := NewElement(line, i, indentTop, nil, tpl, nil)
				if err != nil {
					return nil, err
				}
				tpl.AppendElement(e)
				if err := appendChildren(e, lines, &i, &l, indentTop, e.RawContent, e.Type, tpl); err != nil {
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
	baseDir, err := os.Getwd()
	if err != nil {
		baseDir = ""
	}
	return &Generator{cache: cache, templates: make(map[string]*template.Template), gtemplates: make(map[string]*Template), htmls: make(map[string]string), baseDir: baseDir, delimLeft: defaultDelimLeft, delimRight: defaultDelimRight}
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
func appendChildren(parent Container, lines []string, i *int, l *int, parentIndent int, parentRawContent bool, parentType string, tpl *Template) error {
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
				if err := appendChild(parent, &line, &indent, lines, i, l, tpl); err != nil {
					return err
				}
			}
		default:
			switch {
			case indent < parentIndent+1:
				return nil
			case indent == parentIndent+1:
				if err := appendChild(parent, &line, &indent, lines, i, l, tpl); err != nil {
					return err
				}
			case indent > parentIndent+1:
				return fmt.Errorf("the indent of the line %d is invalid. [template: %s][lineno: %d][line: %s]", *i+1, tpl.Path, *i+1, strings.TrimSpace(line))
			}
		}
	}
	return nil
}

// appendChild appends the child element to the parent element.
func appendChild(parent Container, line *string, indent *int, lines []string, i *int, l *int, tpl *Template) error {
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
	err = appendChildren(child, lines, i, l, child.Indent, child.RawContent, child.Type, tpl)
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
