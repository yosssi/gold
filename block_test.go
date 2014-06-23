package gold

import (
	"bytes"
	"testing"
)

func TestBlockAppendChild(t *testing.T) {
	b := &Block{}
	e := &Element{}
	b.AppendChild(e)
	if len(b.Elements) != 1 || b.Elements[0] != e {
		t.Error("An element was not set to the block.")
	}
}

func TestBlockHtml(t *testing.T) {
	g := NewGenerator(false)
	tpl := NewTemplate("/", g)
	b := &Block{}
	b.Template = tpl
	e, err := NewElement("div#id.class attr=val This is a text.", 1, 0, nil, tpl, b)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	b.AppendChild(e)
	var bf bytes.Buffer
	b.Html(&bf, nil)
	if bf.String() != `<div id="id" class="class" attr="val">This is a text.</div>` {
		t.Errorf("Html returns an invalid string.")
	}
}
