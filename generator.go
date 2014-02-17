package gold

// A generator represents an HTML generator.
type generator struct {
	cache bool
}

// Html parses a template and returns an html string.
func (g *generator) Html(path string, data interface{}) (string, error) {
	return "", nil
}

// NewGenerator generages a generator and returns it.
func NewGenerator(cache bool) generator {
	return generator{cache: cache}
}
