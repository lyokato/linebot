package linebot

import (
	"bytes"
	"encoding/json"
)

func unmarshalJSON(data []byte, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}
