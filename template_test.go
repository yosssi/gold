package gold

import (
	"testing"
)

func TestTemplateAppendElement(t *testing.T) {
	tpl := &Template{}
	e := &Element{}
	tpl.AppendElement(e)
	if len(tpl.Elements) != 1 || tpl.Elements[0] != e {
		t.Errorf("The template's elements are invalid.")
	}
}

func TestTemplateHtml(t *testing.T) {
	super := &Template{}
	e, err := NewElement("div#id.class attr=val This is a text.", 1, 0, nil, super, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	super.AppendElement(e)
	tpl := &Template{Super: super}
	s, err := tpl.Html(nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if s != `<div id="id" class="class" attr="val">This is a text.</div>` {
		t.Errorf("A string is invalid.")
	}
}

func TestTemplateDir(t *testing.T) {
	// When Path has no "/"s.
	tpl := Template{Path: "test.gold"}
	if tpl.Dir() != "./" {
		t.Errorf("Dir return value is invalid.")
	}

	// When Path has a "/".
	tpl = Template{Path: "./test.gold"}
	if tpl.Dir() != "./" {
		t.Errorf("Dir return value is invalid.")
	}

	// When Path has "/"s.
	tpl = Template{Path: "./views/test.gold"}
	if tpl.Dir() != "./views/" {
		t.Errorf("Dir return value is invalid.")
	}
}

func TestTemplateAddBlock(t *testing.T) {
	tpl := &Template{Blocks: map[string]*Block{}}
	b := &Block{Name: "name"}
	tpl.AddBlock(b.Name, b)
	if len(tpl.Blocks) != 1 || tpl.Blocks[b.Name] != b {
		t.Errorf("The template's blocks are invalid.")
	}
}

func TestNewTemplate(t *testing.T) {
	path := "./test.gold"
	g := &Generator{}
	tpl := NewTemplate(path, g)
	if tpl.Path != path || tpl.Generator != g {
		t.Errorf("The template is invalid.")
	}
}
