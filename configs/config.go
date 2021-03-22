package configs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// ETCDConfig represents the configuration information for a etcd cluster
type ETCDConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

// Config represents the configuration information for cb-network
type Config struct {
	ETCD ETCDConfig `yaml:"etcd_cluster"`
}

// LoadConfig represents a function to read a MQTT Broker's configuration information from a file
func LoadConfig(path string) (Config, error) {

	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config, err
}
