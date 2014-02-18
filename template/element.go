package template

import (
	"errors"
	"fmt"
	"strings"
)

// An element represents an html element.
type Element struct {
	text       string
	tokens     []string
	lineNo     int
	indent     int
	parent     *Element
	children   []*Element
	tag        string
	attributes map[string]string
	id         string
	classes    []string
	textValues string
}

// parse parses the element.
func (e *Element) parse() error {
	if e.hasNoTokens() {
		return errors.New(fmt.Sprintf("The element has no tokens. (line no: %d)", e.lineNo))
	}
	for i, token := range e.tokens {
		switch {
		case i == 0:
			if err := e.setTag(token); err != nil {
				return err
			}
		case e.hasTextValues():
		case attribute(token):
		default:
		}
	}
	return nil
}

// hasNoTokens returns if the element has no tokens.
func (e *Element) hasNoTokens() bool {
	return len(e.tokens) == 0
}

// hasTextValues returns if the element has textValues.
func (e *Element) hasTextValues() bool {
	return len(e.textValues) > 0
}

// setTag extracts a tag from the token and sets it to the element.
func (e *Element) setTag(token string) error {
	tag := strings.Split(strings.Split(token, "#")[0], ".")[0]
	if tag == "" {
		tag = "div"
	}
	e.tag = tag
	err := e.setIdFromToken(token)
	if err != nil {
		return err
	}
	e.appendClassesFromToken(token)
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
		e.multipleIdsError()
	}
	return nil
}

// setId sets the id to the element.
func (e *Element) setId(id string) error {
	if e.id != "" {
		return e.multipleIdsError()
	}
	e.id = id
	return nil
}

// multipleIdsError returns a multiple ids error.
func (e *Element) multipleIdsError() error {
	return errors.New(fmt.Sprintf("The number of the element id has to be one. (line no: %d)", e.lineNo))
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

// appendClass appends the class to the element.
func (e *Element) appendClass(class string) {
	e.classes = append(e.classes, class)
}

// NewElement generates a new element and returns it.
func NewElement(text string, lineNo int, indent int, parent *Element) (Element, error) {
	text = strings.TrimSpace(text)
	tokens := tokens(text)
	e := Element{text: text, tokens: tokens, lineNo: lineNo, indent: indent, parent: parent}
	err := e.parse()
	if err != nil {
		return Element{}, err
	}
	return e, nil
}

// tokens returns the string's tokens.
func tokens(s string) []string {
	tokens := make([]string, 0)
	var joinedTokens []string
	joined := false
	for _, token := range strings.Split(s, " ") {
		if joined {
			joinedTokens = append(joinedTokens, token)
			if closed(token) {
				tokens = append(tokens, strings.Join(joinedTokens, " "))
				joinedTokens = make([]string, 0)
				joined = false
			}
		} else {
			if unclosed(token) {
				joined = true
				joinedTokens = []string{token}
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
func unclosed(token string) bool {
	return len(strings.Split(token, "\"")) == 2
}

// unclosed returns if the token is closed or not.
func closed(token string) bool {
	return strings.HasSuffix(token, "\"")
}

// attribute returns if the token is a attribute set or not.
func attribute(token string) bool {
	return strings.Index(token, "=") >= 0
}
