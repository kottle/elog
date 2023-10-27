package filter

import (
	"easylog/internal/common"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ConditionArr struct {
	If    []string `yaml:"if"`
	NotIf []string `yaml:"notIf"`
}
type ConditionMap struct {
	If    map[string]string `yaml:"if"`
	NotIf map[string]string `yaml:"notIf"`
}

type Data struct {
	Field ConditionArr `yaml:"field"`
	Line  ConditionMap `yaml:"line"`
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
	if len(filterData.Field.If) > 0 {
		skipField = true
		for _, include := range filterData.Field.If {
			if include == k {
				skipField = false
			}
		}
	}
	if len(filterData.Field.NotIf) > 0 {
		skipField = false
		for _, exclude := range filterData.Field.NotIf {
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
	f.lock.RLock()
	defer f.lock.RUnlock()
	filterData := f.filterData

	if len(filterData.Line.NotIf) > 0 {
		for k, v := range filterData.Line.NotIf {
			if kvs[k] == v {
				logrus.Debugf("Skip line: with %s:%s", k, v)
				return true
			}
		}
		return false
	}

	if len(filterData.Line.If) > 0 {
		for k, v := range filterData.Line.If {
			if kvs[k] != v {
				logrus.Debugf("Not skip line: with %s:%s", k, v)
				return true
			}
		}
		return false
	}
	return false
}
