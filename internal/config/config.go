package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Представление конфига
type Config struct {
	SourceUrl string `yaml:"source_url"`
	DbFile    string `yaml:"db_file"`
}

// Загрузка конфига
func LoadConfig(path string) (Config, error) {
	// читаем конфиг
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to load config file, err: %w", err)
	}

	// преобразуем в структуру
	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("yaml decode error, err: %w", err)
	}

	return config, nil
}
