package genutil

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var once sync.Once
var instance *configManager

type configManager struct {
	conf ConfigType
}

// Read and get config.yaml
func ConfigFile() *configManager {
	once.Do(func() {
		instance = &configManager{}
	})
	return instance
}

func (s *configManager) Get() ConfigType {
	return s.conf
}

func (s *configManager) Read(filePath string) error {
	ConfigType, err := readConfigFile(filePath)
	if err != nil {
		return err
	}
	s.conf = *ConfigType

	return nil
}

func readConfigFile(filePath string) (*ConfigType, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var ConfigType ConfigType

	err = yaml.Unmarshal(yamlFile, &ConfigType)
	if err != nil {
		return nil, err
	}

	return &ConfigType, nil
}
