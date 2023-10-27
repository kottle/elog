package filter

import (
	"easylog/internal/common"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Field struct {
	Includes []string `yaml:"includes"`
	Excludes []string `yaml:"excludes"`
}
type FieldValue struct {
	Includes common.KVS `yaml:"includes"`
	Excludes common.KVS `yaml:"excludes"`
}
type Data struct {
	Field      Field      `yaml:"field"`
	FieldValue FieldValue `yaml:"field_value"`
}

func New(filepath string) (*Filter, error) {
	f := &Filter{
		filepath: filepath,
	}
	err := f.init()
	return f, err
}

type Filter struct {
	filepath   string
	lock       sync.RWMutex
	filterData Data
}

func (f *Filter) init() error {
	yfile, err := os.ReadFile(f.filepath)
	if err != nil {
		return err
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	err = yaml.Unmarshal(yfile, &f.filterData)
	if err != nil {
		return err
	}
	return nil
}

// SkipField filters the key value pairs
// retuns true if the field should be skipped
func (f *Filter) SkipField(k string) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	var skipField bool
	filterData := f.filterData
	if len(filterData.Field.Includes) > 0 {
		skipField = true
		for _, include := range filterData.Field.Includes {
			if include == k {
				skipField = false
			}
		}
	}
	if len(filterData.Field.Excludes) > 0 {
		skipField = false
		for _, exclude := range filterData.Field.Excludes {
			if exclude == k {
				skipField = true
			}
		}
	}
	return skipField
}

// SkipLine filters the key value pairs
// retuns true if the line should be skipped
func (f *Filter) SkipLine(kvs common.KVS) bool {
	return false
}
