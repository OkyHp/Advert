package main

import (
	"fmt"
	"os"

	"github.com/untrustedmodders/go-plugify"
	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	Host   string `yaml:"host"`
	Base   string `yaml:"base"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Port   uint16 `yaml:"port"`
	Schema string `yaml:"schema"`
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
	pathToFile := fmt.Sprintf("%s/advert.yml", plugify.ConfigsDir)

	file, err := os.Open(pathToFile)
	if err != nil {
		return ConfigData{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config ConfigData

	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return ConfigData{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, err
}
