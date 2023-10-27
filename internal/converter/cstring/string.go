package cstring

import (
	"easylog/internal/color"
	"easylog/internal/common"
)

func Convert(kvs common.KVS) string {
	result := ""
	if file, ok := kvs["file"]; ok {
		result += color.Red(file) + " "
	}
	if time, ok := kvs["time"]; ok {
		result += color.Yellow(time) + " "
	}
	if level, ok := kvs["level"]; ok {
		result += color.Yellow(level) + " "
	}
	if msg, ok := kvs["msg"]; ok {
		result += color.Red(msg) + " "
	}
	for k, v := range kvs {
		if k != "time" && k != "level" && k != "msg" && k != "file" {
			result += k + "=" + v + " "
		}
	}
	return result
}
