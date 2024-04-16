package app

import (
	"log"

	"github.com/ponchik327/Yadro-Golang/tree/main/internal/config"
	"github.com/ponchik327/Yadro-Golang/tree/main/internal/utilsDb"
)

// Логика работы приложения
func RunApp() {
	// парсим флаги
	pathConfig := utilsDb.ParseFlags()

	// загружаем конфиг из config.yaml
	config, err := config.LoadConfig(pathConfig)
	if err != nil {
		log.Fatal("error open config: " + err.Error())
	}

	// создаём бд
	db, err := utilsDb.CreateDatabase(config.SourceUrl, config.DbFile, config.NumGorutine)
	if err != nil {
		log.Fatal("error create database: " + err.Error())
	}
	defer db.Close()
}
