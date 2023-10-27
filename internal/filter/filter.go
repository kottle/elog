package filter

import (
	"easylog/internal/common"
	"os"
	"regexp"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

type ConditionArr struct {
	If    []string `yaml:"if"`
	NotIf []string `yaml:"notIf"`
}
type ConditionMap struct {
	If       map[string][]string `yaml:"if"`
	NotIf    map[string][]string `yaml:"notIf"`
	NotMatch map[string][]string `yaml:"notMatch"`
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
	watcher    *fsnotify.Watcher
}

func (f *Filter) init() error {
	err := f.updateFile()
	if err != nil {
		return err
	}
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = f.watcher.Add(f.filepath)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event, ok := <-f.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.Debugf("modified file: %s", event.Name)
					err := f.updateFile()
					if err != nil {
						logrus.Error(err)
					}
				}
			case err, ok := <-f.watcher.Errors:
				if !ok {
					return
				}
				logrus.Error(err)
			}
		}
	}()

	return nil
}

func (f *Filter) updateFile() error {
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
func (f *Filter) Close() error {
	return f.watcher.Close()
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
		skipLine := false
		for k, v := range filterData.Line.NotIf {
			if slices.Contains(v, kvs[k]) {
				logrus.Debugf("Skip line: with %s:%s", k, v)
				skipLine = true
				break
			}
		}
		for k, regexs := range filterData.Line.NotMatch {
			for _, regex := range regexs {
				if kvs[k] != "" {
					matched, _ := regexp.MatchString(regex, kvs[k])
					if matched {
						logrus.Debugf("Skip line: with %s:%s", k, regex)
						skipLine = true
						break
					}
				}
			}
		}
		return skipLine
	}

	if len(filterData.Line.If) > 0 {
		for k, v := range filterData.Line.If {
			if !slices.Contains(v, kvs[k]) {
				logrus.Debugf("Not skip line: with %s:%s", k, v)
				return true
			}
		}
		return false
	}
	return false
}
