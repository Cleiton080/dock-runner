package main

import (
	"github.com/Cleiton080/dock-runner/internal/configwatcher"
	"github.com/Cleiton080/dock-runner/pkg/config"
)

func main() {
	configChannel := make(chan *config.Config)

	defer close(configChannel)

	configWatcher := configwatcher.NewConfigWatcher("./config.toml")

	if err := configWatcher.Watch(configChannel); err != nil {
		panic(err)
	}

	for conf := range configChannel {
		// enviar config para o control plane
		println(conf)
	}

}
