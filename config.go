package main

import (
	"fmt"
	"os"

	"github.com/untrustedmodders/go-plugify"
	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	Host   string            `yaml:"host"`
	Base   string            `yaml:"base"`
	User   string            `yaml:"user"`
	Pass   string            `yaml:"pass"`
	Port   uint16            `yaml:"port"`
	Params map[string]string `yaml:"params"`
}

type ConfigData struct {
	Debug           bool     `yaml:"debug"`
	TimerInterval   float64  `yaml:"timerInterval"`
	ServerId        uint16   `yaml:"serverId"`
	ServerIp        string   `yaml:"serverIp"`
	HtmlMsgDuration int32    `yaml:"htmlMsgDuration"`
	Database        DBConfig `yaml:"database"`
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

	MSGDebug("Advert ReadConfig: %v", config)

	return config, err
}
