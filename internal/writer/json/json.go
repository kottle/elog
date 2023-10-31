package json

import (
	"easylog/internal/common"
	ejson "encoding/json"
)

type Writer struct {
}

func New(filepath string) *Writer {
	return &Writer{}
}
func (c *Writer) Write(kvs common.KVS) string {
	return write(kvs)
}
func write(kvs common.KVS) string {
	//marshal to json
	b, err := ejson.Marshal(kvs)
	if err == nil {
		return string(b)
	}
	return ""

}
