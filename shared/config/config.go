package config

import (
	"encoding/json"
	"errors"
	"os"
)

type ConfigType string

const (
	Ntfy    ConfigType = "ntfy"
	Discord ConfigType = "discord"
)

var initialized = false

var configs Configs

type NotificationConfig struct {
	Type                 ConfigType `json:"type"`
	Enabled              bool       `json:"enabled"`
	Channel              string     `json:"channel"`
	SubscribedTopicsList []string   `json:"topics"`
	AuthUserEnv          string     `json:"auth_user_env"`
	AuthPassEnv          string     `json:"auth_pass_env"`
	AuthTokenEnv         string     `json:"auth_token_env"`
	SubscribedTopics     map[string]any
}

type Configs struct {
	NotificationConfigs []NotificationConfig `json:"notification_configs"`
}

func GetConfigs() Configs {
	loadConfigs()
	return configs
}

func ReadConfigFile() error {
	initialized = false
	return loadConfigs()
}

func loadConfigs() error {
	if initialized {
		return nil
	}

	configFile, err := os.Open("config.json")
	if err != nil {
		return errors.New("Failed to open config file: " + err.Error())
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&configs)
	if err != nil {
		return errors.New("Malformed config file: " + err.Error())
	}
	initialized = true
	initializeTopicsSet()
	return nil
}

func initializeTopicsSet() {
	for _, config := range configs.NotificationConfigs {
		config.SubscribedTopics = make(map[string]any)
		for _, topic := range config.SubscribedTopicsList {
			config.SubscribedTopics[topic] = nil
		}
	}
}
