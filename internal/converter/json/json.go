package json

import (
	"easylog/internal/common"
	ejson "encoding/json"
)

func Convert(kvs common.KVS) string {
	//marshal to json
	b, err := ejson.Marshal(kvs)
	if err == nil {
		return string(b)
	}
	return ""

}
