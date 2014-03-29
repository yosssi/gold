package gold

import (
	"html/template"
	"testing"
)

func TestGeneratorParseFile(t *testing.T) {
	// When cashe is true and a cached template exists.
	tmplt := &template.Template{}
	g := &Generator{cache: true, templates: map[string]*template.Template{"path": tmplt}}
	tpl, err := g.ParseFile("path")
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if tpl != tmplt {
		t.Errorf("Returned value is invalid.")
	}

	// When g.Parse returns an error.
	g = &Generator{}
	_, err = g.ParseFile("./test/TestGeneratorParseFile/001.gold")
	expectedErrMsg := "The indent of the line 2 is invalid."
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When gtpl.Html() returns an error.
	g = &Generator{}
	_, err = g.ParseFile("./test/TestGeneratorParseFile/002.gold")
	expectedErrMsg = "The block element does not have a name. (line no: 1)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When teplate.Parse() returns an error.
	g = &Generator{}
	_, err = g.ParseFile("./test/TestGeneratorParseFile/003.gold")
	expectedErrMsg = "template: ./test/TestGeneratorParseFile/003.gold:1: missing value for command"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When cache is true and g.ParseFile returns tpl.
	g = NewGenerator(true)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/004.gold")
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
}
func TestGeneratorParseString(t *testing.T) {
	g := &Generator{}
	parent := `
doctype html
html
  head
    title Gold
  body
    block content
    footer
      block footer
`
	child := `
extends parent

block content
  #container
    | Hello Gold

block footer
  .footer
    | Copyright XXX
	include inc
`
	inc := `
p This is an included line.
`
	stringTemplates := map[string]string{"parent": parent, "child": child, "inc": inc}
	_, err := g.ParseString(stringTemplates, "child")
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
}

func TestParse(t *testing.T) {
	// When cache is true and Parse returns a cached template.
	gtmplt := &Template{}
	g := NewGenerator(true)
	g.gtemplates = map[string]*Template{"path": gtmplt}
	gtpl, err := g.parse("path", nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	if gtmplt != gtpl {
		t.Errorf("Returned value is invalid.")
	}

	// When ioutil.ReadFile returns an error.
	gtmplt = &Template{}
	g = NewGenerator(false)
	gtpl, err = g.parse("./somepath/somefile", nil)
	expectedErrMsg := "open ./somepath/somefile: no such file or directory"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When a template includes an empty line.
	g = NewGenerator(false)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/005.gold")
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// When a template includes a "extends" and returns an error.
	g = NewGenerator(false)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/006.gold")
	expectedErrMsg = "The line tokens length is invalid. (expected: 2, actual: 1, line no: 1)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When a template includes a "extends" and returns an error while parsing a super template.
	g = NewGenerator(false)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/007.gold")
	expectedErrMsg = "open ./test/TestGeneratorParseFile/./somepath/somefile.gold: no such file or directory"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When a template includes a "extends" and returns an error while parsing a block line.
	g = NewGenerator(false)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/008.gold")
	expectedErrMsg = "The lien tokens length is invalid. (expected: 2, actual: 1, line no: 3)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When a template includes a "extends" and returns an error while appending a child.
	g = NewGenerator(false)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/009.gold")
	expectedErrMsg = "The indent of the line 4 is invalid."
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When NewElement returns an error.
	g = NewGenerator(false)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/010.gold")
	expectedErrMsg = "The number of the element id has to be one. (line no: 1)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When NewElement returns an error.
	g = NewGenerator(false)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/011.gold")
	expectedErrMsg = "The number of the element id has to be one. (line no: 4)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When cache is true
	g = NewGenerator(true)
	_, err = g.ParseFile("./test/TestGeneratorParseFile/012.gold")
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
}

func TestNewGenerator(t *testing.T) {
	// When cache is true.
	g := NewGenerator(true)
	if g.cache != true {
		t.Errorf("g.cache should be true.")
	}

	// When cache is false.
	g = NewGenerator(false)
	if g.cache != false {
		t.Errorf("g.cache should be false.")
	}
}

func TestFormatLf(t *testing.T) {
	if formatLf("a\n\r\r\n") != "a\n\n\n" {
		t.Errorf("Return value is invalid.")
	}
}

func TestIndent(t *testing.T) {
	// When a string includes unicodeTab.
	if indent(string(unicodeTab)+"a") != 1 {
		t.Errorf("Return value is invalid.")
	}

	// When a string includes double unicodeSpaces.
	if indent(string(unicodeSpace)+string(unicodeSpace)+"a") != 1 {
		t.Errorf("Return value is invalid.")
	}

	// When a string includes no unicodeTabs and no unicodeSpaces.
	if indent("a") != 0 {
		t.Errorf("Return value is invalid.")
	}
}

func TestEmpty(t *testing.T) {
	// When the stirng is "".
	if empty("") != true {
		t.Errorf("Return value is invalid.")
	}

	// When the stirng is unicodeTab.
	if empty(string(unicodeTab)) != true {
		t.Errorf("Return value is invalid.")
	}

	// When the stirng is unicodeSpace.
	if empty(string(unicodeSpace)) != true {
		t.Errorf("Return value is invalid.")
	}
}

func TestTopElement(t *testing.T) {
	// When the string is a top element.
	if topElement("a") != true {
		t.Errorf("Return value is invalid.")
	}

	// When the string is not a top element.
	if topElement("  a") != false {
		t.Errorf("Return value is invalid.")
	}
}

func TestAppendChildren(t *testing.T) {
	// When line is empty.
	i := 0
	l := 1
	err := appendChildren(nil, []string{""}, &i, &l, 0, false, "")
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// parentRawContent is true and indent < parentIndent+1
	i = 0
	l = 1
	err = appendChildren(nil, []string{"a"}, &i, &l, 0, true, "")
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// parentRawContent is true and indent >= parentIndent+1 and appendChild returns an error.
	i = 0
	l = 1
	e, err := NewElement("", 0, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	err = appendChildren(e, []string{"  div#id1#id2"}, &i, &l, 0, true, "")
	expectedErrMsg := "The number of the element id has to be one. (line no: 1)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// parentType is TypeBlock and indent < parentIndent+1
	i = 0
	l = 1
	err = appendChildren(nil, []string{"a"}, &i, &l, 0, false, TypeBlock)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// parentType is TypeBlock and indent >= parentIndent+1
	i = 0
	l = 1
	err = appendChildren(nil, []string{"  a"}, &i, &l, 0, false, TypeBlock)
	expectedErrMsg = "The indent of the line 1 is invalid. Block element can not have child elements."
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// parentType is TypeTag and indent < parentIndent+1
	i = 0
	l = 1
	err = appendChildren(nil, []string{"a"}, &i, &l, 0, false, TypeTag)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// parentType is TypeTag and indent == parentIndent+1 and appendChild returns an error.
	e, err = NewElement("", 0, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	i = 0
	l = 1
	err = appendChildren(e, []string{"  div#id1#id2"}, &i, &l, 0, false, TypeTag)
	expectedErrMsg = "The number of the element id has to be one. (line no: 1)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// parentType is TypeTag and indent > parentIndent+1.
	e, err = NewElement("", 0, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	i = 0
	l = 1
	err = appendChildren(e, []string{"    div"}, &i, &l, 0, false, TypeTag)
	expectedErrMsg = "The indent of the line 1 is invalid."
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// parentType is TypeTag and indent == parentIndent+1 and appendChild returns no errors.
	e, err = NewElement("", 0, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	i = 0
	l = 1
	err = appendChildren(e, []string{"  div"}, &i, &l, 0, false, TypeTag)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
}

func TestAppendChild(t *testing.T) {
	// When parent is *Block.
	block := &Block{}
	i := 0
	l := 1
	line := "div"
	indent := 0
	err := appendChild(block, &line, &indent, []string{"div"}, &i, &l)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// When parent is *Element.
	e, err := NewElement("", 0, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	i = 0
	l = 1
	line = "div"
	indent = 0
	err = appendChild(e, &line, &indent, []string{"div"}, &i, &l)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// When parent is *Element and NewElement returns an error.
	e, err = NewElement("", 0, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	i = 0
	l = 1
	line = "div#id1#id2"
	indent = 0
	err = appendChild(e, &line, &indent, []string{"div#id1#id2"}, &i, &l)
	expectedErrMsg := "The number of the element id has to be one. (line no: 1)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}

	// When parent is *Element and appendChildren returns an error.
	e, err = NewElement("", 0, 0, nil, nil, nil)
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}
	i = 0
	l = 2
	line = "div"
	indent = 0
	err = appendChild(e, &line, &indent, []string{"div", "  div#id3#id4"}, &i, &l)
	expectedErrMsg = "The number of the element id has to be one. (line no: 2)"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}
}

func TestIsExtends(t *testing.T) {
	// When the line has "extends".
	if isExtends("extends aaa") != true {
		t.Errorf("Return value is invalid.")
	}

	// When the line does not have "extends".
	if isExtends("aaa") != false {
		t.Errorf("Return value is invalid.")
	}
}

func TestIsBlock(t *testing.T) {
	// When the line has "block".
	if isBlock("block aaa") != true {
		t.Errorf("Return value is invalid.")
	}

	// When the line does not have "block".
	if isBlock("aaa") != false {
		t.Errorf("Return value is invalid.")
	}
}
