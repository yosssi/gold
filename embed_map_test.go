package gold

import (
	"testing"
)

func TestNewEmbedMap(t *testing.T) {
	// When no error occurs.
	_, err := NewEmbedMap([]string{"name=Foo"})
	if err != nil {
		t.Errorf("An error(%s) occurred.", err.Error())
	}

	// When an error occurs.
	_, err = NewEmbedMap([]string{"name"})
	expectedErrMsg := "the parameter did not have = and a key-value could not be derived. [parameter: name]"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("An error(%s) should be occurred.", expectedErrMsg)
	}
}
