package gold

import (
	"bytes"
)

// A Block represents a Block of a Gold template.
type Block struct {
	Name     string
	Elements []*Element
	Template *Template
}

// AppendChild appends the element to the receiver block.
func (b *Block) AppendChild(child *Element) {
	b.Elements = append(b.Elements, child)
}

// Html writes the block's html to the buffer.
func (b *Block) Html(bf *bytes.Buffer) {
	for _, e := range b.Elements {
		e.Html(bf)
	}
}
