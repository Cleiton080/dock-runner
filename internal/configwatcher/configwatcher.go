package configwatcher

import (
	"github.com/fsnotify/fsnotify"

	"github.com/Cleiton080/dock-runner/pkg/config"
)

type ConfigWatcher struct {
	FilePath string
}

func NewConfigWatcher(filePath string) *ConfigWatcher {
	return &ConfigWatcher{
		FilePath: filePath,
	}
}

func (cw *ConfigWatcher) Watch(configChannel chan *config.Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := watcher.Add(cw.FilePath); err != nil {
		return err
	}

	configFile := config.NewConfigFile(cw.FilePath)

	for event := range watcher.Events {
		if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
			config, err := configFile.Load()
			if err != nil {
				return err
			}

			configChannel <- config
		}
	}

	return nil
}
