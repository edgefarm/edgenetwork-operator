package json

import (
	"bytes"
	"encoding/json"
)

func Marshal(t interface{}, escapeHTML bool) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(escapeHTML)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
