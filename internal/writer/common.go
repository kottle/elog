package writer

import "easylog/internal/common"

type IWriter interface {
	Write(common.KVS) string
}
