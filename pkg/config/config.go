package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type ConfigFile struct {
	Path string
}

type Config struct {
	Containers []Container
}

type PlatformType string

const (
	PlatformLinuxArm64 PlatformType = "linux/arm64"
	PlatformLinuxAmd64 PlatformType = "linux/amd64"
)

type Container struct {
	Token      string
	Count      int32
	Platform   PlatformType
	Ephemeral  bool
	Repository ContainerRepository
	Volumes    []ContainerVolume
}

type ContainerRepository struct {
	Name  string
	Owner string
}

type ContainerVolume struct {
	HostPath      string
	ContainerPath string
}

func NewConfigFile(configPath string) *ConfigFile {
	return &ConfigFile{
		Path: configPath,
	}
}

func (cf ConfigFile) Load() (*Config, error) {
	var config Config

	if _, err := os.Stat(cf.Path); os.IsNotExist(err) {
		return nil, fmt.Errorf("the config file '%s' doesn't exist", cf.Path)
	}

	if file, err := os.ReadFile(cf.Path); err == nil {
		if err := toml.Unmarshal(file, &config); err != nil {
			return nil, err
		}

		return &config, nil
	}

	return nil, fmt.Errorf("the config file '%s' could not be read, please check its readable permission", cf.Path)
}
