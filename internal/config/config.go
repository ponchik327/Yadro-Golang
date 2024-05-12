package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Представление конфига
type Config struct {
	SourceUrl   string `yaml:"source_url"`
	DbFile      string `yaml:"db_file"`
	IndexFile   string `yaml:"index_file"`
	NumGorutine int    `yaml:"parallel"`
}

// Загрузка конфига
func LoadConfig(path string) (Config, error) {
	// читаем конфиг
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	// преобразуем в структуру
	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
