package gold

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
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

func TestHasNoTokens(t *testing.T) {
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

func TestHasTextValues(t *testing.T) {
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

func TestParseFirstToken(t *testing.T) {
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

func TestSetTag(t *testing.T) {
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

func TestSetIdFromToken(t *testing.T) {
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

func TestSetId(t *testing.T) {
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

func TestMultipleIdsError(t *testing.T) {
	e := &Element{Attributes: make(map[string]string), LineNo: 1}
	expectedErrMsg := fmt.Sprintf("The number of the element id has to be one. (line no: %d)", e.LineNo)
	if err := e.multipleIdsError(); err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Error(%s) should be returned.", expectedErrMsg)
	}
}
