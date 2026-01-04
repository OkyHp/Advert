package main

import (
	"fmt"
	"os"

	"github.com/untrustedmodders/go-plugify"
	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	Host   string            `yml:"host"`
	Base   string            `yml:"base"`
	User   string            `yml:"user"`
	Pass   string            `yml:"pass"`
	Port   uint16            `yml:"port"`
	Params map[string]string `yml:"params"`
}

type ConfigData struct {
	TimerInterval   float64  `yml:"timerInterval"`
	ServerId        uint16   `yml:"serverId"`
	ServerIp        string   `yml:"serverIp"`
	HtmlMsgDuration int32    `yml:"htmlMsgDuration"`
	Database        DBConfig `yml:"database"`
}

func ReadConfig() (ConfigData, error) {
	pathToFile := fmt.Sprintf("%s/config.yml", plugify.Plugin.Location)

	data, err := os.ReadFile(pathToFile)
	if err != nil {
		return ConfigData{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ConfigData

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return ConfigData{}, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return config, err
}
