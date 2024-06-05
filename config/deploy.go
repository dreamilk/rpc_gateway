package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type DeployConfig struct {
	Addr   string `yaml:"addr"`
	Consul string `yaml:"consul"`
}

const (
	fileName = "./deploy.yaml"
)

var (
	DeployConf DeployConfig
)

func init() {
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("init deploy.yaml error %s", err)
	}

	if err := yaml.Unmarshal(b, &DeployConf); err != nil {
		log.Fatalf("init deploy.yaml error %s", err)
	}
}
