package main

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	ConfigFilePaths = []string{
		".notepad.yml",
		path.Join(getHomeDir(), ".notepad.yml"),
		path.Join(getHomeDir(), "go-notepad", "notepad.yml"),
	}
)

var DefaultConfig = ConfigSchema{
	Font: ConfigFont{
		Family: "Lucida Console",
		Size:   10,
	},
	StatusBar: ConfigStatusBar{
		Enable: false,
	},
}

type (
	ConfigSchema struct {
		Font      ConfigFont
		StatusBar ConfigStatusBar
	}

	ConfigFont struct {
		Family string
		Size   int64
	}

	ConfigStatusBar struct {
		Enable bool
	}
)

func loadConfig(filePath string) (*ConfigSchema, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML file into the Config struct
	var config ConfigSchema
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func searchAndLoadConfig() (*ConfigSchema, error) {
	for _, c := range ConfigFilePaths {
		if fileExist(c) {
			return loadConfig(c)
		}
	}

	return &DefaultConfig, nil
}
