package app

import (
	"fmt"
	"log"
	"time"

	"github.com/ponchik327/Yadro-Golang/tree/main/internal/config"
	"github.com/ponchik327/Yadro-Golang/tree/main/internal/search"
	"github.com/ponchik327/Yadro-Golang/tree/main/internal/utilsDb"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/index"
)

// Логика работы приложения
func RunApp() {
	// парсим флаги
	pathConfig, searchLine, isIndexSearch := utilsDb.ParseFlags()

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

	// создаём индекс
	index, err := index.CreateIndex(config.IndexFile, db)
	if err != nil {
		log.Fatal("error create index: " + err.Error())
	}

	start := time.Now()

	search.Search(searchLine, db, index, isIndexSearch)

	duration := time.Since(start)
	fmt.Println("time search: ")
	fmt.Println(duration.Microseconds())
}
