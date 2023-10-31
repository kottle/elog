package cstring

import (
	"easylog/internal/color"
	"easylog/internal/common"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Style struct {
	Foreground string `yaml:"fg"`
	Background string `yaml:"bg"`
}
type ThemeData struct {
	Styles map[string]Style `yaml:"styles"`
}

type Writer struct {
	filepath string
	lock     sync.RWMutex
	theme    ThemeData
	watcher  *fsnotify.Watcher
}

func New(filepath string) (*Writer, error) {
	w := &Writer{
		filepath: filepath,
	}
	err := w.init()
	return w, err
}

func (w *Writer) init() error {
	err := w.updateFile()
	if err != nil {
		return err
	}
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = w.watcher.Add(w.filepath)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.Debugf("modified file: %s", event.Name)
					err := w.updateFile()
					if err != nil {
						logrus.Error(err)
					}
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				logrus.Error(err)
			}
		}
	}()

	return nil
}

func (w *Writer) updateFile() error {
	yfile, err := os.ReadFile(w.filepath)
	if err != nil {
		return err
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	err = yaml.Unmarshal(yfile, &w.theme)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) Close() error {
	return w.watcher.Close()
}

func (c *Writer) Write(kvs common.KVS) string {
	return write(kvs)
}

func write(kvs common.KVS) string {
	result := ""
	if file, ok := kvs["@TAG"]; ok {
		result += color.Red(file) + " "
	}
	if time, ok := kvs["time"]; ok {
		result += color.Yellow(time) + " "
	}
	if msg, ok := kvs["@time"]; ok {
		result += color.Yellow(msg) + " "
	}
	if level, ok := kvs["level"]; ok {
		result += color.Yellow(level) + " "
	}
	if level, ok := kvs["@level"]; ok {
		result += color.Yellow(level) + " "
	}
	if msg, ok := kvs["msg"]; ok {
		result += color.Red(msg) + " "
	}
	if msg, ok := kvs["@message"]; ok {
		result += color.Red(msg) + " "
	}
	for k, v := range kvs {
		if k != "time" && k != "level" && k != "msg" && k != "file" && k != "@message" && k != "@time" && k != "@level" && k != "@TAG" {
			result += k + "=" + v + " "
		}
	}
	return result
}
