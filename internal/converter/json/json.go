package json

import (
	ejson "encoding/json"
)

func Convert(kvs map[string]string) string {
	//marshal to json
	b, err := ejson.Marshal(kvs)
	if err == nil {
		return string(b)
	}
	return ""

}
