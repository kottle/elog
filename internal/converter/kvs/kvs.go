package kvs

import (
	"strings"
)

func ToKVS(line string, filter func(k, v string) bool) map[string]string {
	//parse line to json
	line = strings.ReplaceAll(line, "\\\"", "'")
	line = strings.ReplaceAll(line, "\\\\'", "'")
	kvs := make(map[string]string)
	remaining := line
	for {
		var key string
		var value string
		//logrus.Debugf("remaining: %s", remaining)
		pos := strings.Index(remaining, "=")
		if pos > 0 {
			key = strings.Trim(remaining[:pos], " ")
			//logrus.Debug("key: ", key)
			remaining = remaining[pos+1:]
			if strings.Index(remaining, "=") <= 0 {
				if remaining[0] == '"' {
					remaining = remaining[1 : len(remaining)-1]
				}
				value = remaining
				//logrus.Debug("value: ", value)
				remaining = ""

			} else if remaining[0] == '"' {
				pos = strings.Index(remaining[1:], "\" ")
				if pos > 0 {
					value = remaining[1 : pos+1]
					//logrus.Debug("value: ", value)
					remaining = remaining[pos+2:]
				}
			} else {
				pos = strings.Index(remaining, " ")
				if pos > 0 {
					value = remaining[:pos]
					//logrus.Debug("value: ", value)
					remaining = remaining[pos+1:]
				}
			}
			if filter != nil && !filter(key, value) {
				continue
			}
			kvs[key] = value
		} else {
			//end line
			break
		}
	}
	return kvs
}
