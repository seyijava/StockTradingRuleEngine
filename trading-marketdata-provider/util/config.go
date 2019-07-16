package util

import (
	"encoding/json"
	"fmt"
	"os"
)

type kafkaConfig struct {
	BrokerHost string `json:"brokerHost"`
	Topic      string `json:"topic"`
	Partition  int32  `json:"partition"`
}

type AppConfig struct {
	KafkaConfig kafkaConfig `json:"kafka"`
}

func LoadConfiguration(file string) (AppConfig, error) {
	var config AppConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, err
}
