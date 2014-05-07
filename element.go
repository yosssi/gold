package gold

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	TypeTag               = "tag"
	TypeContent           = "content"
	TypeBlock             = "block"
	TypeExpression        = "expression"
	TypeLiteral           = "literal"
	TypeInclude           = "include"
	TypeOutputExpression  = "outputExpression"
	IncludeParaStartIndex = 2
)

var (
	doctypes = map[string]string{
		"html":         "<!DOCTYPE html>",
		"xml":          "<?xml version=\"1.0\" encoding=\"utf-8\" ?>",
		"transitional": "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">",
		"strict":       "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\">",
		"frameset":     "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\">",
		"1.1":          "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">",
		"basic":        "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML Basic 1.1//EN\" \"http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd\">",
		"mobile":       "<!DOCTYPE html PUBLIC \"-//WAPFORUM//DTD XHTML Mobile 1.2//EN\" \"http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd\">",
	}
)

// An Element represents an element of a Gold template.
type Element struct {
	Text             string
	Tokens           []string
	LineNo           int
	Indent           int
	Parent           *Element
	Children         []*Element
	Tag              string
	Attributes       map[string]string
	SingleAttributes []string
	Id               string
	Classes          []string
	TextValues       []string
	Type             string
	Template         *Template
	Block            *Block
	RawContent       bool
}

// parse parses the element.
func (e *Element) parse() error {
	if e.hasNoTokens() {
		return errors.New(fmt.Sprintf("The element has no tokens. (line no: %d)", e.LineNo))
	}
	switch {
	case e.Type == TypeContent || e.Type == TypeBlock || e.Type == TypeExpression || e.Type == TypeLiteral || e.Type == TypeInclude || e.Type == TypeOutputExpression || e.comment():
	default:
		for i, token := range e.Tokens {
			switch {
			case i == 0:
				if err := e.parseFirstToken(token); err != nil {
					return err
				}
			case e.hasTextValues():
				e.appendTextValue(token)
			case attribute(token):
				e.appendAttribute(token)
			case singleAttribute(token):
				e.appendSingleAttribute(token)
			default:
				e.appendTextValue(token)
			}
		}
	}
	return nil
}

// hasNoTokens returns if the element has no tokens.
func (e *Element) hasNoTokens() bool {
	return len(e.Tokens) == 0
}

// hasTextValues returns if the element has textValues.
func (e *Element) hasTextValues() bool {
	return len(e.TextValues) > 0
}

// parseFirstToken parses the token and sets values to the element.
func (e *Element) parseFirstToken(token string) error {
	switch token {
	case "javascript:":
		e.setTag("script")
		e.appendAttribute("type=text/javascript")
	default:
		if err := e.setTag(token); err != nil {
			return err
		}
	}
	return nil
}

// setTag extracts a tag from the token and sets it to the element.
func (e *Element) setTag(token string) error {
	tag := strings.Split(strings.Split(token, "#")[0], ".")[0]
	if tag == "" {
		tag = "div"
	}
	e.Tag = tag
	err := e.setIdFromToken(token)
	if err != nil {
		return err
	}
	e.appendClassesFromToken(token)
	if e.Tag == "script" || e.Tag == "style" || strings.HasSuffix(token, ".") {
		e.RawContent = true
	}
	return nil
}

// setIdFromToken extracts an id from the token and sets it to the element.
func (e *Element) setIdFromToken(token string) error {
	parts := strings.Split(token, "#")
	switch len(parts) {
	case 1:
	case 2:
		e.setId(strings.Split(parts[1], ".")[0])
	default:
		return e.multipleIdsError()
	}
	return nil
}

// setId sets the id to the element.
func (e *Element) setId(id string) error {
	if e.Id != "" {
		return e.multipleIdsError()
	}
	e.Id = id
	return nil
}

// multipleIdsError returns a multiple ids error.
func (e *Element) multipleIdsError() error {
	return errors.New(fmt.Sprintf("The number of the element id has to be one. (line no: %d)", e.LineNo))
}

// appendClassesFromToken extracts classes from the token and appends them to the element.
func (e *Element) appendClassesFromToken(token string) {
	for i, part := range strings.Split(token, ".") {
		if i == 0 {
			continue
		}
		e.appendClass(strings.Split(part, "#")[0])
	}
}

// appendClass appends the class to the element's classes.
func (e *Element) appendClass(class string) {
	if class != "" {
		e.Classes = append(e.Classes, class)
	}
}

// appendTextValue appends the token to the element's textValues.
func (e *Element) appendTextValue(token string) {
	e.TextValues = append(e.TextValues, token)
}

// appendAttribute appends the token to the element's attributes.
func (e *Element) appendAttribute(token string) {
	kv := strings.Split(token, "=")
	if len(kv) < 2 {
		return
	}
	k := kv[0]
	v := parseValue(strings.Join(kv[1:], "="))
	switch k {
	case "id":
		e.setId(v)
	case "class":
		e.appendClass(v)
	default:
		e.Attributes[k] = v
	}
}

func (e *Element) appendSingleAttribute(token string) {
	e.SingleAttributes = append(e.SingleAttributes, strings.TrimSuffix(strings.TrimPrefix(token, "["), "]"))

}

// AppendChild appends the element to the receiver element.
func (e *Element) AppendChild(child *Element) {
	e.Children = append(e.Children, child)
}

// Html writes the element's html to the buffer.
func (e *Element) Html(bf *bytes.Buffer, stringTemplates map[string]string) error {
	switch {
	case e.comment():
	case e.Type == TypeContent || e.Type == TypeExpression || e.Type == TypeOutputExpression:
		if e.Type == TypeOutputExpression {
			bf.WriteString("{{" + strings.Join(e.Tokens[1:], " ") + "}}")
		} else {
			e.writeText(bf)
		}
		for _, child := range e.Children {
			err := child.Html(bf, stringTemplates)
			if err != nil {
				return err
			}
		}
	case e.Type == TypeLiteral:
		e.writeLiteralValue(bf)
	case e.Type == TypeBlock:
		if len(e.Tokens) < 2 {
			return errors.New(fmt.Sprintf("The block element does not have a name. (line no: %d)", e.LineNo))
		}
		name := e.Tokens[1]
		sub := e.getTemplate().Sub
		if sub == nil {
			if err := e.writeChildren(bf, stringTemplates); err != nil {
				return err
			}
			return nil
		}
		block := sub.Blocks[name]
		if block == nil {
			if err := e.writeChildren(bf, stringTemplates); err != nil {
				return err
			}
			return nil
		}
		block.Html(bf, stringTemplates)
	case e.Type == TypeInclude:
		if len(e.Tokens) < 2 {
			return errors.New(fmt.Sprintf("The include element does not have a path. (line no: %d)", e.LineNo))
		}
		tpl := e.getTemplate()
		incTplPath := e.Tokens[1]
		g := tpl.Generator
		var incTpl *Template
		var err error
		addBaseDir := true
		if stringTemplates == nil {
			if g.baseDir != "" && CurrentDirectoryBasedPath(incTplPath) {
				incTplPath = tpl.Dir() + incTplPath
				addBaseDir = false
			}
			incTpl, err = g.parse(incTplPath+Extension, nil, addBaseDir)
		} else {
			incTpl, err = g.parse(incTplPath, stringTemplates, addBaseDir)
		}
		if err != nil {
			return err
		}
		embedMap, err := NewEmbedMap(e.Tokens[IncludeParaStartIndex:])
		if err != nil {
			return err
		}
		incHtml, err := incTpl.Html(stringTemplates, embedMap)
		if err != nil {
			return err
		}
		bf.WriteString(incHtml)
	default:
		e.writeOpenTag(bf)
		if e.hasTextValues() {
			e.writeTextValue(bf)
		}
		if err := e.writeChildren(bf, stringTemplates); err != nil {
			return err
		}
		e.writeCloseTag(bf)
	}
	return nil
}

// writeChildren writes the element's children's HTML.
func (e *Element) writeChildren(bf *bytes.Buffer, stringTemplates map[string]string) error {
	for _, child := range e.Children {
		err := child.Html(bf, stringTemplates)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeOpenTag writes the element's open tag to the buffer.
func (e *Element) writeOpenTag(bf *bytes.Buffer) {
	switch e.Tag {
	case "doctype":
		if doctype, prs := doctypes[e.textValue()]; prs {
			bf.WriteString(doctype)
		} else {
			bf.WriteString("<!DOCTYPE ")
			bf.WriteString(e.textValue())
			bf.WriteString(">")
		}
	default:
		bf.WriteString("<")
		bf.WriteString(e.Tag)
		if e.hasId() {
			e.writeId(bf)
		}
		if e.hasClasses() {
			e.writeClasses(bf)
		}
		if e.hasAttributes() {
			e.writeAttributes(bf)
		}
		if e.hasSingleAttributes() {
			e.writeSingleAttributes(bf)
		}
		bf.WriteString(">")
	}
}

// writeText writes the element's text to the buffer.
func (e *Element) writeText(bf *bytes.Buffer) {
	bf.WriteString(e.Text + "\n")
}

// textValue returns the element's textValues.
func (e *Element) textValue() string {
	return strings.Join(e.TextValues, " ")
}

// hasId returns if the element has an id or not.
func (e *Element) hasId() bool {
	return e.Id != ""
}

// writeId writes the element's id to the buffer.
func (e *Element) writeId(bf *bytes.Buffer) {
	bf.WriteString(" id=\"")
	bf.WriteString(e.Id)
	bf.WriteString("\"")
}

// hasClasses returns if the element has classes or not.
func (e *Element) hasClasses() bool {
	return len(e.Classes) > 0
}

// writeClasses writes the element's classes to the buffer.
func (e *Element) writeClasses(bf *bytes.Buffer) {
	bf.WriteString(" class=\"")
	for i, class := range e.Classes {
		if i > 0 {
			bf.WriteString(" ")
		}
		bf.WriteString(class)
	}
	bf.WriteString("\"")
}

// hasAttributes returns if the element has attributes or not.
func (e *Element) hasAttributes() bool {
	return len(e.Attributes) > 0
}

// hasSingleAttributes returns if the element has single attributes or not.
func (e *Element) hasSingleAttributes() bool {
	return len(e.SingleAttributes) > 0
}

// writeAttributes writes the element's attributes to the buffer.
func (e *Element) writeAttributes(bf *bytes.Buffer) {
	for k, v := range e.Attributes {
		bf.WriteString(" ")
		bf.WriteString(k)
		bf.WriteString("=\"")
		bf.WriteString(v)
		bf.WriteString("\"")
	}
}

// writeSingleAttributes writes the element's single attributes to the buffer.
func (e *Element) writeSingleAttributes(bf *bytes.Buffer) {
	for _, v := range e.SingleAttributes {
		bf.WriteString(" ")
		bf.WriteString(v)
	}
}

// writeTextValue writes the element's text value to the buffer.
func (e *Element) writeTextValue(bf *bytes.Buffer) {
	switch e.Tag {
	case "doctype":
	default:
		bf.WriteString(e.textValue())
	}
}

// writeCloseTag writes the element's close tag to the buffer.
func (e *Element) writeCloseTag(bf *bytes.Buffer) {
	switch e.Tag {
	case "doctype":
	default:
		bf.WriteString("</")
		bf.WriteString(e.Tag)
		bf.WriteString(">")
	}
}

// setType sets a type to the element.
func (e *Element) setType() {
	switch {
	case e.Parent != nil && (e.Parent.RawContent || e.Parent.Type == TypeContent):
		e.Type = TypeContent
	case len(e.Tokens) > 0 && e.Tokens[0] == "block":
		e.Type = TypeBlock
	case len(e.Tokens) > 0 && e.Tokens[0] == "include":
		e.Type = TypeInclude
	case len(e.Tokens) > 0 && e.Tokens[0] == "|":
		e.Type = TypeLiteral
	case expression(e.Text):
		e.Type = TypeExpression
	case len(e.Tokens) > 0 && e.Tokens[0] == "=":
		e.Type = TypeOutputExpression
	default:
		e.Type = TypeTag
	}
}

// getTemplate returns the element's template.
func (e *Element) getTemplate() *Template {
	switch {
	case e.Parent != nil:
		return e.Parent.getTemplate()
	case e.Block != nil:
		return e.Block.Template
	default:
		return e.Template
	}
}

// literalValue returns the element's literal value.
func (e *Element) literalValue() string {
	if len(e.Tokens) < 2 {
		return ""
	}
	return strings.Join(e.Tokens[1:], " ")
}

// writeLiteralValue writes the element's literal value to the buffer.
func (e *Element) writeLiteralValue(bf *bytes.Buffer) {
	bf.WriteString(e.literalValue())
}

// comment returns if the string is a comment or not.
func (e *Element) comment() bool {
	return strings.HasPrefix(e.Text, "//")
}

// NewElement generates a new element and returns it.
func NewElement(text string, lineNo int, indent int, parent *Element, tpl *Template, block *Block) (*Element, error) {
	rawText := text
	text = strings.TrimSpace(text)
	tokens := tokens(text)
	e := &Element{Text: text, Tokens: tokens, LineNo: lineNo, Indent: indent, Parent: parent, Attributes: make(map[string]string), Template: tpl, Block: block}
	e.setType()
	if e.Type == TypeContent {
		e.Text = rawText
	}
	err := e.parse()
	if err != nil {
		return nil, err
	}
	return e, nil
}

// tokens returns the string's tokens.
func tokens(s string) []string {
	tokens := make([]string, 0)
	var joinedTokens []string
	joined := false
	var closeMark string
	for _, token := range strings.Split(s, " ") {
		if joined {
			joinedTokens = append(joinedTokens, token)
			if closed(token, closeMark) {
				tokens = append(tokens, strings.Join(joinedTokens, " "))
				joinedTokens = make([]string, 0)
				joined = false
			}
		} else {
			if unclosed, mark := unclosed(token); unclosed {
				joined = true
				joinedTokens = []string{token}
				closeMark = mark
			} else {
				tokens = append(tokens, token)
			}
		}
	}
	if joined {
		tokens = append(tokens, joinedTokens...)
	}
	return tokens
}

// unclosed returns if the token is unclosed or not.
func unclosed(token string) (bool, string) {
	if len(strings.Split(token, "\"")) == 2 {
		return true, `"`
	}
	if strings.HasPrefix(token, "[") && !strings.HasSuffix(token, "]") {
		return true, "]"
	}

	return false, ""
}

// unclosed returns if the token is closed or not.
func closed(token string, closeMark string) bool {
	return strings.HasSuffix(token, closeMark)
}

// attribute returns if the token is a attribute set or not.
func attribute(token string) bool {
	return strings.Index(token, "=") >= 0
}

// singleAttribute returns if the token is a single attribute set or not.
func singleAttribute(token string) bool {
	return strings.HasPrefix(token, "[") && strings.HasSuffix(token, "]")
}

// parseValue parses the value and returns a result string.
func parseValue(value string) string {
	if literal(value) {
		return value[1 : len(value)-1]
	}
	return value
}

// literal returns the string is a literal or not.
func literal(s string) bool {
	return len(s) > 1 && s[0] == unicodeDoubleQuote && s[len(s)-1] == unicodeDoubleQuote
}

// expression returns the string is an expression or not.
func expression(s string) bool {
	return strings.HasPrefix(s, "{{") && strings.HasSuffix(s, "}}")
}
