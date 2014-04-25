package gold

import (
	"bytes"
	"fmt"
	"testing"
)

func TestElementParse(t *testing.T) {
	// When an element has no tokens.
	e := &Element{LineNo: 5}
	err := e.parse()
	expectedErrMsg := fmt.Sprintf("The element has no tokens. (line no: %d)", e.LineNo)
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When an element's type is TypeContent.
	e = &Element{Tokens: []string{"test"}, Type: TypeContent}
	err = e.parse()
	if err != nil {
		t.Errorf("Error(%s) occurred.", err.Error())
	}

	// When an element's type is TypeBlock.
	e = &Element{Tokens: []string{"test"}, Type: TypeBlock}
	err = e.parse()
	if err != nil {
		t.Errorf("Error(%s) occurred.", err.Error())
	}

	// When an element's type is TypeExpression.
	e = &Element{Tokens: []string{"test"}, Type: TypeExpression}
	err = e.parse()
	if err != nil {
		t.Errorf("Error(%s) occurred.", err.Error())
	}

	// When an element's type is TypeLiteral.
	e = &Element{Tokens: []string{"test"}, Type: TypeLiteral}
	err = e.parse()
	if err != nil {
		t.Errorf("Error(%s) occurred.", err.Error())
	}

	// When an element's type is TypeTag.
	e, err = NewElement("div data-test=test test test2", 1, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("Error(%s) occurred.", err.Error())
	}
	err = e.parse()
	if err != nil {
		t.Errorf("Error(%s) occurred.", err.Error())
	}
}

func TestElementHasNoTokens(t *testing.T) {
	// When an element has no tokens.
	e := &Element{Tokens: nil}
	if e.hasNoTokens() != true {
		t.Errorf("hasNoTokens should return true.")
	}

	// When an element has tokens.
	e = &Element{Tokens: []string{"test", "test2"}}
	if e.hasNoTokens() != false {
		t.Errorf("hasNoTokens should return false.")
	}
}

func TestElementHasTextValues(t *testing.T) {
	// When an element has text values.
	e := &Element{TextValues: []string{"test"}}
	if e.hasTextValues() != true {
		t.Errorf("hasTextValues should return true.")
	}

	// When an element has no text values.
	e = &Element{TextValues: nil}
	if e.hasTextValues() != false {
		t.Errorf("hasTextValues should return false.")
	}
}

func TestElementParseFirstToken(t *testing.T) {
	// When parsing "javascript:".
	e := &Element{Attributes: make(map[string]string)}
	if err := e.parseFirstToken("javascript:"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e.Tag != "script" {
		t.Errorf(`Tag "script" was not set to the element.`)
	}
	if a := e.Attributes; len(a) < 1 || a["type"] != "text/javascript" {
		t.Errorf(`Attribute "type=text/javascript" was not set to the element.`)
	}
	if e.RawContent != true {
		t.Errorf("RawContent should be true.")
	}

	// When parsing a token which has multiple ids.
	e = &Element{Attributes: make(map[string]string)}
	if err := e.parseFirstToken("div#id1#id2"); err == nil {
		t.Errorf("No errors occurred.")
	}

	// When parsing a token.
	e = &Element{Attributes: make(map[string]string)}
	if err := e.parseFirstToken("div#id.class"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e.Tag != "div" {
		t.Errorf(`Tag "div" was not set to the element.`)
	}
	if e.Id != "id" {
		t.Errorf(`Id should be "id".`)
	}
	if len(e.Classes) != 1 || e.Classes[0] != "class" {
		t.Errorf(`Id should be ["class"].`)
	}
}

func TestElementSetTag(t *testing.T) {
	// When a tag is "".
	e := &Element{Attributes: make(map[string]string)}
	if err := e.setTag(".class"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// When a token has multiple ids.
	e = &Element{Attributes: make(map[string]string)}
	if err := e.setTag("div#id1#id2"); err == nil {
		t.Errorf("No errors occurred.")
	}

	// When a tag is "script".
	e = &Element{Attributes: make(map[string]string)}
	if err := e.setTag("script"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e.RawContent != true {
		t.Errorf("RawContent should be true.")
	}

	// When a tag is "style".
	e = &Element{Attributes: make(map[string]string)}
	if err := e.setTag("style"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e.RawContent != true {
		t.Errorf("RawContent should be true.")
	}

	// When a tag is "p.".
	e = &Element{Attributes: make(map[string]string)}
	if err := e.setTag("style"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e.RawContent != true {
		t.Errorf("RawContent should be true.")
	}
}

func TestElementSetIdFromToken(t *testing.T) {
	// When a token has no ids.
	e := &Element{Attributes: make(map[string]string)}
	if err := e.setIdFromToken("div"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e.Id != "" {
		t.Errorf("Id should be empty.")
	}

	// When a token has an id.
	e = &Element{Attributes: make(map[string]string)}
	if err := e.setIdFromToken("div#id"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e.Id != "id" {
		t.Errorf(`id should be "id".`)
	}

	// When a token has multiple ids.
	e = &Element{Attributes: make(map[string]string), LineNo: 1}
	expectedErrMsg := fmt.Sprintf("The number of the element id has to be one. (line no: %d)", e.LineNo)
	if err := e.setIdFromToken("div#id1#id2"); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}
}

func TestElementSetId(t *testing.T) {
	// When setting an id to the element.
	e := &Element{Attributes: make(map[string]string), LineNo: 1}
	if err := e.setId("id"); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e.Id != "id" {
		t.Errorf(`id should be "id".`)
	}

	// When setting multiple ids to the element.
	expectedErrMsg := fmt.Sprintf("The number of the element id has to be one. (line no: %d)", e.LineNo)
	if err := e.setId("id2"); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}
}

func TestElementMultipleIdsError(t *testing.T) {
	e := &Element{Attributes: make(map[string]string), LineNo: 1}
	expectedErrMsg := fmt.Sprintf("The number of the element id has to be one. (line no: %d)", e.LineNo)
	if err := e.multipleIdsError(); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}
}

func TestElementAppendClassesFromToken(t *testing.T) {
	e := &Element{Attributes: make(map[string]string)}
	e.appendClassesFromToken("div.class1#id.class2")
	if len(e.Classes) != 2 || e.Classes[0] != "class1" || e.Classes[1] != "class2" {
		t.Errorf("The element's classes are invalid.")
	}
}

func TestElementAppendClass(t *testing.T) {
	// A class is "".
	e := &Element{Attributes: make(map[string]string)}
	e.appendClass("")
	if len(e.Classes) != 0 {
		t.Errorf("The element's classes are invalid.")
	}

	// A class is not "".
	e = &Element{Attributes: make(map[string]string)}
	e.appendClass("test")
	if len(e.Classes) != 1 || e.Classes[0] != "test" {
		t.Errorf("The element's classes are invalid.")
	}
}

func TestElementAppendTextValue(t *testing.T) {
	e := &Element{Attributes: make(map[string]string)}
	e.appendTextValue("test")
	if len(e.TextValues) != 1 || e.TextValues[0] != "test" {
		t.Errorf("The element's text values are invalid.")
	}
}

func TestElementAppendAttribute(t *testing.T) {
	// When a token has no "=".
	e := &Element{Attributes: make(map[string]string)}
	e.appendAttribute("test")
	if len(e.Attributes) != 0 {
		t.Errorf("The element's attributes are invalid.")
	}

	// When a token has a "id=".
	e = &Element{Attributes: make(map[string]string)}
	e.appendAttribute("id=testid")
	if e.Id != "testid" {
		t.Errorf("The element's id is invalid.")
	}

	// When a token has a "class=".
	e = &Element{Attributes: make(map[string]string)}
	e.appendAttribute("class=testclass")
	if len(e.Classes) != 1 || e.Classes[0] != "testclass" {
		t.Errorf("The element's class is invalid.")
	}

	// When a token has an attribute.
	e = &Element{Attributes: make(map[string]string)}
	e.appendAttribute("data-test=testdata")
	if len(e.Attributes) != 1 || e.Attributes["data-test"] != "testdata" {
		t.Errorf("The element's attributes are invalid.")
	}

	// When a token has multiple  "="s.
	e = &Element{Attributes: make(map[string]string)}
	e.appendAttribute("data-test=testdata=testdata")
	if len(e.Attributes) != 1 || e.Attributes["data-test"] != "testdata=testdata" {
		t.Errorf("The element's attributes are invalid.")
	}
}

func TestElementAppendChild(t *testing.T) {
	e := &Element{Attributes: make(map[string]string)}
	child := &Element{Attributes: make(map[string]string)}
	e.AppendChild(child)
	if len(e.Children) != 1 || e.Children[0] != child {
		t.Errorf("The element's chilredn are invalid.")
	}
}

func TestElementHtml(t *testing.T) {
	// When the element's type is expression and child.Html returns an error.
	parent, err := NewElement("{{}}", 1, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	child, err := NewElement("block", 2, 1, parent, nil, nil)
	parent.AppendChild(child)
	var bf bytes.Buffer
	expectedErrMsg := fmt.Sprintf("The block element does not have a name. (line no: %d)", child.LineNo)
	if err := parent.Html(&bf, nil); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When the element's type is literal.
	e, err := NewElement("| abc", 1, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	bf = bytes.Buffer{}
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if bf.String() != "abc" {
		t.Errorf("Html output is invalid.")
	}

	// When the element's type is block and the template's sub is nil.
	tpl := &Template{}
	e, err = NewElement("block test", 1, 0, nil, tpl, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	bf = bytes.Buffer{}
	expectedErrMsg = "The template does not have a sub template."
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// When the element's type is block and the template's sub's block is nil.
	tpl = &Template{Sub: &Template{Blocks: make(map[string]*Block)}}
	e, err = NewElement("block test", 1, 0, nil, tpl, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	bf = bytes.Buffer{}
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// When the element's type is block.
	block := &Block{Name: "test"}
	blockElement, err := NewElement("div#id.class attr=val This is a text.", 1, 0, nil, nil, block)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	block.AppendChild(blockElement)
	tpl = &Template{Sub: &Template{Blocks: map[string]*Block{"test": block}}}
	e, err = NewElement("block test", 1, 0, nil, tpl, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	bf = bytes.Buffer{}
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	expectedString := `<div id="id" class="class" attr="val">This is a text.</div>`
	if bf.String() != expectedString {
		t.Errorf("Buffer stirng should be %s", expectedString)
	}

	// When the element's type is tag and child.Html returns an error.
	e, err = NewElement("p This is a text.", 1, 0, nil, nil, nil)
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	child, err = NewElement("block", 2, 1, e, nil, nil)
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	e.AppendChild(child)
	bf = bytes.Buffer{}
	expectedErrMsg = fmt.Sprintf("The block element does not have a name. (line no: %d)", child.LineNo)
	if err := e.Html(&bf, nil); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When the element's type is tag and child.Html returns an error.
	e, err = NewElement("div.class", 1, 0, nil, nil, nil)
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	child, err = NewElement("p This is a text.", 2, 1, e, nil, nil)
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	e.AppendChild(child)
	bf = bytes.Buffer{}
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	expectedString = `<div class="class"><p>This is a text.</p></div>`
	if bf.String() != expectedString {
		t.Errorf("Buffer stirng should be %s", expectedString)
	}

	// When the element's type is include and tokens' length < 2.
	e, err = NewElement("include", 1, 0, nil, nil, nil)
	expectedErrMsg = fmt.Sprintf("The include element does not have a path. (line no: %d)", e.LineNo)
	if err := e.Html(&bf, nil); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When the element's type is include and g.Parse returns an error.
	g := NewGenerator(false)
	tpl = NewTemplate("path", g)
	e, err = NewElement("include ./somepath/somefile", 1, 0, nil, tpl, nil)
	bf = bytes.Buffer{}
	expectedErrMsg = "open ././somepath/somefile.gold: no such file or directory"
	if err := e.Html(&bf, nil); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When the element's type is include and incTpl.Html returns an error.
	g = NewGenerator(false)
	tpl = NewTemplate("./test/TestElementHtml/somefile.gold", g)
	e, err = NewElement("include ./001", 1, 0, nil, tpl, nil)
	bf = bytes.Buffer{}
	expectedErrMsg = "The block element does not have a name. (line no: 1)"
	if err := e.Html(&bf, nil); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When the element's type is include.
	g = NewGenerator(false)
	tpl = NewTemplate("./test/TestElementHtml/somefile.gold", g)
	e, err = NewElement("include ./002", 1, 0, nil, tpl, nil)
	bf = bytes.Buffer{}
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// When the element's type is include.
	g = NewGenerator(false)
	tpl = NewTemplate("./test/TestElementHtml/somefile.gold", g)
	e, err = NewElement("include ./002 param", 1, 0, nil, tpl, nil)
	bf = bytes.Buffer{}
	expectedErrMsg = "the parameter did not have = and a key-value could not be derived. [parameter: param]"
	if err := e.Html(&bf, nil); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When the block's sub template does not exist and an error occurs.
	g = NewGenerator(false)
	tpl = NewTemplate("./test/TestElementHtml/003.gold", g)
	parent, err = NewElement("block test", 1, 0, nil, tpl, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	child, err = NewElement("block", 2, 1, parent, tpl, nil)
	parent.AppendChild(child)
	bf = bytes.Buffer{}
	expectedErrMsg = "The block element does not have a name. (line no: 2)"
	if err := parent.Html(&bf, nil); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When the block's sub template exists and an error occurs.
	g = NewGenerator(false)
	parentTpl := NewTemplate("./test/TestElementHtml/003.gold", g)
	subTpl := NewTemplate("./test/TestElementHtml/003.gold", g)
	parentTpl.Sub = subTpl
	parent, err = NewElement("block test", 1, 0, nil, parentTpl, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	child, err = NewElement("block", 2, 1, parent, parentTpl, nil)
	parent.AppendChild(child)
	bf = bytes.Buffer{}
	expectedErrMsg = "The block element does not have a name. (line no: 2)"
	if err := parent.Html(&bf, nil); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When the element's type is OutputExpression.
	e, err = NewElement(`= "abc"`, 1, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	bf = bytes.Buffer{}
	if err := e.Html(&bf, nil); err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if bf.String() != `{{"abc"}}` {
		t.Errorf("Html output is invalid. [output: %s]", bf.String())
	}
}

func TestElementWriteOpenTag(t *testing.T) {
	// When element's tag is doctype and a text value is html.
	e, err := NewElement("doctype html", 1, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	bf := bytes.Buffer{}
	e.writeOpenTag(&bf)
	expectedString := `<!DOCTYPE html>`
	if bf.String() != expectedString {
		t.Errorf("Buffer stirng should be %s", expectedString)
	}

	// When element's tag is doctype and the element has a custom text value.
	e, err = NewElement("doctype AABBCC", 1, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	bf = bytes.Buffer{}
	e.writeOpenTag(&bf)
	expectedString = `<!DOCTYPE AABBCC>`
	if bf.String() != expectedString {
		t.Errorf("Buffer stirng should be %s", expectedString)
	}

	// When element's tag is doctype and the element has a custom text value.
	e, err = NewElement("div#id.class attr=val AAABBBCCC", 1, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	bf = bytes.Buffer{}
	e.writeOpenTag(&bf)
	expectedString = `<div id="id" class="class" attr="val">`
	if bf.String() != expectedString {
		t.Errorf("Buffer stirng should be %s", expectedString)
	}
}

func TestElementWriteText(t *testing.T) {
	e := &Element{Text: "This is a text."}
	bf := bytes.Buffer{}
	e.writeText(&bf)
	expectedString := "This is a text.\n"
	if bf.String() != expectedString {
		t.Errorf("Buffer stirng should be %s", expectedString)
	}
}

func TestElementTextValue(t *testing.T) {
	e := &Element{TextValues: []string{"a", "b", "c"}}
	expectedString := "a b c"
	if s := e.textValue(); s != expectedString {
		t.Errorf("Returned stirng should be %s", expectedString)
	}
}

func TestElementHasId(t *testing.T) {
	// When the element's id is empty.
	e := &Element{Id: ""}
	if e.hasId() != false {
		t.Errorf("hasId sholud return false.")
	}

	// When the element's id is not empty.
	e = &Element{Id: "id"}
	if e.hasId() != true {
		t.Errorf("hasId sholud return true.")
	}
}

func TestElementWriteId(t *testing.T) {
	e := &Element{Id: "id"}
	var bf bytes.Buffer
	e.writeId(&bf)
	expectedString := ` id="id"`
	if bf.String() != expectedString {
		t.Errorf("Return string should be %s", expectedString)
	}
}

func TestElementHasClasses(t *testing.T) {
	// When the element has no classes.
	e := &Element{}
	if e.hasClasses() != false {
		t.Errorf("Return value should be false.")
	}

	// When the element has classes.
	e = &Element{Classes: []string{"a"}}
	if e.hasClasses() != true {
		t.Errorf("Return value should be false.")
	}
}

func TestElementWriteClasses(t *testing.T) {
	e := &Element{Classes: []string{"a", "b"}}
	var bf bytes.Buffer
	e.writeClasses(&bf)
	expectedString := ` class="a b"`
	if bf.String() != expectedString {
		t.Errorf("Return string should be %s", expectedString)
	}
}

func TestElementHasAttributes(t *testing.T) {
	// When the element has no attributes.
	e := &Element{}
	if e.hasAttributes() != false {
		t.Errorf("Return value should be false.")
	}

	// When the element has classes.
	e = &Element{Attributes: map[string]string{"a": "b"}}
	if e.hasAttributes() != true {
		t.Errorf("Return value should be false.")
	}
}

func TestElementWriteAttributes(t *testing.T) {
	e := &Element{Attributes: map[string]string{"a": "b"}}
	var bf bytes.Buffer
	e.writeAttributes(&bf)
	expectedString := ` a="b"`
	if bf.String() != expectedString {
		t.Errorf("Return string should be %s", expectedString)
	}
}

func TestElementWriteTextValue(t *testing.T) {
	// When the element's tag is doctype.
	e := &Element{Tag: "doctype", TextValues: []string{"a", "b"}}
	var bf bytes.Buffer
	e.writeTextValue(&bf)
	expectedString := ``
	if bf.String() != expectedString {
		t.Errorf("Return string should be %s", expectedString)
	}

	// When the element's tag is not doctype.
	e = &Element{Tag: "div", TextValues: []string{"a", "b"}}
	bf = bytes.Buffer{}
	e.writeTextValue(&bf)
	expectedString = `a b`
	if bf.String() != expectedString {
		t.Errorf("Return string should be %s", expectedString)
	}
}

func TestElementCloseTag(t *testing.T) {
	// When the element's tag is doctype.
	e := &Element{Tag: "doctype"}
	var bf bytes.Buffer
	e.writeCloseTag(&bf)
	expectedString := ``
	if bf.String() != expectedString {
		t.Errorf("Return string should be %s", expectedString)
	}

	// When the element's tag is not doctype.
	e = &Element{Tag: "div"}
	bf = bytes.Buffer{}
	e.writeCloseTag(&bf)
	expectedString = `</div>`
	if bf.String() != expectedString {
		t.Errorf("Return string should be %s", expectedString)
	}
}

func TestElementSetType(t *testing.T) {
	// When the element's Parent.RawContent is true.
	parent := &Element{RawContent: true}
	e := &Element{Parent: parent}
	e.setType()
	if e.Type != TypeContent {
		t.Errorf("Type should be %s", TypeContent)
	}

	// When the element's Parent.Type is TypeContent.
	parent = &Element{Type: TypeContent}
	e = &Element{Parent: parent}
	e.setType()
	if e.Type != TypeContent {
		t.Errorf("Type should be %s", TypeContent)
	}

	// When the element's first token is block.
	e = &Element{Tokens: []string{"block"}}
	e.setType()
	if e.Type != TypeBlock {
		t.Errorf("Type should be %s", TypeBlock)
	}

	// When the element's first token is |.
	e = &Element{Tokens: []string{"|"}}
	e.setType()
	if e.Type != TypeLiteral {
		t.Errorf("Type should be %s", TypeLiteral)
	}

	// When the element's text is an expression.
	e = &Element{Text: "{{.}}"}
	e.setType()
	if e.Type != TypeExpression {
		t.Errorf("Type should be %s", TypeExpression)
	}

	// When the element's text is an tag element.
	e = &Element{Text: "div"}
	e.setType()
	if e.Type != TypeTag {
		t.Errorf("Type should be %s", TypeTag)
	}

	// When the element's text is an include element.
	e = &Element{Tokens: []string{"include"}}
	e.setType()
	if e.Type != TypeInclude {
		t.Errorf("Type should be %s", TypeInclude)
	}

	// When the element's first token is =.
	e = &Element{Tokens: []string{"="}}
	e.setType()
	if e.Type != TypeOutputExpression {
		t.Errorf("Type should be %s", TypeOutputExpression)
	}

}

func TestElementGetTemplate(t *testing.T) {
	// When the element has a parent.
	template := &Template{}
	parent := &Element{Template: template}
	e := &Element{Parent: parent}
	if tpl := e.getTemplate(); tpl != template {
		t.Errorf("Returned template is invalid.")
	}

	// When the element has a block.
	template = &Template{}
	block := &Block{Template: template}
	e = &Element{Block: block}
	if tpl := e.getTemplate(); tpl != template {
		t.Errorf("Returned template is invalid.")
	}

	// When the element has a template.
	template = &Template{}
	e = &Element{Template: template}
	if tpl := e.getTemplate(); tpl != template {
		t.Errorf("Returned template is invalid.")
	}
}

func TestElementLiteralValue(t *testing.T) {
	// When the element's tokens' length is less than 2.
	e := &Element{}
	if e.literalValue() != "" {
		t.Errorf("Returned value is invalid.")
	}
	// When the element's tokens' length is greater than or equal to 2.
	e = &Element{Tokens: []string{"|", "a", "b"}}
	if e.literalValue() != "a b" {
		t.Errorf("Returned value is invalid.")
	}
}

func TestElementWriteLiteralValue(t *testing.T) {
	e := &Element{Tokens: []string{"|", "a", "b"}}
	var bf bytes.Buffer
	e.writeLiteralValue(&bf)
	if bf.String() != "a b" {
		t.Errorf("Return string is invalid.")
	}
}

func TestElementComment(t *testing.T) {
	// When the element is a comment.
	e := &Element{Text: "//aaa"}
	if e.comment() != true {
		t.Errorf("Return value should be true.")
	}

	// When the element is not a comment.
	e = &Element{Text: "aaa"}
	if e.comment() != false {
		t.Errorf("Return value should be true.")
	}
}

func TestNewElement(t *testing.T) {
	// When an error occurs while parsing.
	_, err := NewElement("div#id1#id2", 1, 0, nil, nil, nil)
	expectedErrMsg := fmt.Sprintf("The number of the element id has to be one. (line no: %d)", 1)
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When an element is returned.
	e, err := NewElement("div", 1, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if e == nil || e.Tag != "div" {
		t.Errorf("Returned value is invalid.")
	}
}

func TestTokens(t *testing.T) {
	// When a pair of double quotes exists.
	text := `div attr="val1 val2" AAA`
	tkns := tokens(text)
	if len(tkns) != 3 || tkns[0] != "div" || tkns[1] != `attr="val1 val2"` || tkns[2] != "AAA" {
		t.Errorf("Returned value is invalid.")
	}

	// When a double quote exists.
	text = `div "AAA BBB`
	tkns = tokens(text)
	if len(tkns) != 3 || tkns[0] != "div" || tkns[1] != `"AAA` || tkns[2] != `BBB` {
		t.Errorf("Returned value is invalid.")
	}
}

func TestUnclosed(t *testing.T) {
	// When the token is unclosed.
	if unclosed, _ := unclosed(`aaa"bbb`); unclosed != true {
		t.Errorf("Returned value should be true.")
	}

	// When the token is closed.
	if unclosed, _ := unclosed(`aaa"bbb"`); unclosed != false {
		t.Errorf("Returned value should be false.")
	}
}

func TestClosed(t *testing.T) {
	// When the token is closed.
	if closed(`aaabbb"`, `"`) != true {
		t.Errorf("Returned value should be true.")
	}

	// When the token is unclosed.
	if closed("aaa", `"`) != false {
		t.Errorf("Returned value should be false.")
	}
}

func TestAttribute(t *testing.T) {
	// When the token is an attribute.
	if attribute("a=b") != true {
		t.Errorf("Returned value should be true.")
	}

	// When the token is not an attribute.
	if attribute("abc") != false {
		t.Errorf("Returned value should be false.")
	}
}

func TestParseValue(t *testing.T) {
	// When the value is literal.
	if parseValue(`"aaa"`) != "aaa" {
		t.Errorf("Returned value is invalid.")
	}

	// When the value is not literal.
	if parseValue("aaa") != "aaa" {
		t.Errorf("Returned value is invalid.")
	}
}

func TestLiteral(t *testing.T) {
	// When the value is literal.
	if literal(`"aaa"`) != true {
		t.Errorf("Returned value should be true.")
	}

	// When the value is not literal.
	if literal("aaa") != false {
		t.Errorf("Returned value should be false.")
	}
}

func TestExpression(t *testing.T) {
	// When the value is expression.
	if expression("{{.}}") != true {
		t.Errorf("Returned value should be true.")
	}

	// When the value is not expression.
	if expression("aaa") != false {
		t.Errorf("Returned value should be false.")
	}
}
