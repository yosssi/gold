package gold

import (
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

	// When an element's type is TypeExpression.
	e = &Element{Tokens: []string{"test"}, Type: TypeLiteral}
	err = e.parse()
	if err != nil {
		t.Errorf("Error(%s) occurred.", err.Error())
	}

	// When an element's type is TypeExpression.
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
