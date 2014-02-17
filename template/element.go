package template

// An element represents an html element.
type element struct {
	text       string
	tokens     []string
	lineNo     int
	indent     int
	parent     *element
	children   []*element
	tag        string
	attributes map[string]string
	id         string
	classes    []string
}
