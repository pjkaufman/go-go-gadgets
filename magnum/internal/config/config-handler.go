package config

import (
	"encoding/json"
	"fmt"

	commandhandler "github.com/pjkaufman/go-go-gadgets/pkg/command-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

const (
	configDirName  = "magnum"
	configFileName = "series.json"
)

func WriteConfig(config *Config) {
	if config == nil {
		return
	}

	configDir := getConfigLocation()
	err := filehandler.MustCreateFolderIfNotExists(configDir)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to json marshal config: %s", err))
	}

	jsonConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to json marshal config: %s", err))
	}

	configFile := filehandler.JoinPath(configDir, configFileName)

	filehandler.WriteFileContents(configFile, string(jsonConfig))
}

func GetConfig() *Config {
	configDir := getConfigLocation()
	folderExists, err := filehandler.FolderExists(configDir)
	if err != nil {
		logger.WriteError(err.Error())
	}

	if !folderExists {
		return &Config{}
	}

	configFile := filehandler.JoinPath(configDir, configFileName)
	fileExists, err := filehandler.FileExists(configFile)
	if err != nil {
		logger.WriteError(err.Error())
	}

	if !fileExists {
		return &Config{}
	}

	jsonConfig := filehandler.ReadInFileContents(configFile)
	var config = &Config{}

	err = json.Unmarshal([]byte(jsonConfig), config)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to json unmarshal config from %q: %s", configFile, err))
	}

	return config
}

func getConfigLocation() string {
	userConfigDir := commandhandler.MustGetUserConfigDir()

	return filehandler.JoinPath(userConfigDir, configDirName)
}
