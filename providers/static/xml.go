package static

import (
	"bytes"
	"encoding/json"

	"github.com/basgys/goxml2json"
)

type xmlDecoder struct {}

func (i *xmlDecoder) Decode(data []byte, v interface{}) error {
	buf, err := xml2json.Convert(bytes.NewReader(data))
	if err != nil {
		return err
	}
	return json.Unmarshal(buf.Bytes(), &v)
}
