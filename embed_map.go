package gold

import (
	"fmt"
	"strings"
)

// An EmbedMap represents a map for embedding strings to Gold tempaltes.
type EmbedMap map[string]string

// A NewEmbedMap generates an EmbedMap and returns it.
func NewEmbedMap(kvs []string) (EmbedMap, error) {
	embedMap := EmbedMap{}
	for _, kv := range kvs {
		kvMap := strings.Split(kv, "=")
		if len(kvMap) != 2 {
			return nil, fmt.Errorf("the parameter did not have = and a key-value could not be derived. [parameter: %s]", kv)
		}
		embedMap[trimDoubleQuote(kvMap[0])] = trimDoubleQuote(kvMap[1])
	}
	return embedMap, nil
}

// trimDoubleQuote trims double quotes.
func trimDoubleQuote(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)
}
