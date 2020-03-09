package static

import (
	"encoding/json"

	"github.com/ghodss/yaml"
)


type yamlDecoder struct {}

func (i *yamlDecoder) Decode(data []byte, v interface{}) error {
	bb, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bb, &v)
}
