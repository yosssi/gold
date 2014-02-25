package gold

// A Container represents a container which holds elements.
type Container interface {
	AppendChild(*Element)
}
