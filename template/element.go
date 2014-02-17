package template

import (
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
}

// NewElement generates a new element and returns it.
func NewElement(text string, lineNo int, indent int, parent *Element) (Element, error) {
	text = strings.TrimSpace(text)
	return Element{}, nil
}

// tokens returns the string's tokens.
func tokens(s string) []string {
	tokens := make([]string, 0)
	var joinTokens []string
	join := false
	for _, token := range strings.Split(s, " ") {
		if join {
			joinTokens = append(joinTokens, token)
			if strings.HasSuffix(token, "\"") {
				tokens = append(tokens, strings.Join(joinTokens, " "))
				joinTokens = make([]string, 0)
				join = false
			}
		} else {
			if unclosed(token) {
				join = true
				joinTokens = []string{token}
			} else {
				tokens = append(tokens, token)
			}
		}
	}
	if join {
		tokens = append(tokens, joinTokens...)
	}
	return tokens
}

func unclosed(s string) bool {
	return len(strings.Split(s, "\"")) == 2
}
