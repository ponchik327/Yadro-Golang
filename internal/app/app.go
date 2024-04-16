package app

import (
	"fmt"
	"log"
	"os"

	"github.com/ponchik327/Yadro-Golang/tree/main/internal/config"
	"github.com/ponchik327/Yadro-Golang/tree/main/internal/utilsDb"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/database"
)

// Логика работы приложения
func RunApp() {
	pathConfig, needShowDb, numComics := utilsDb.ParseFlags()

	// Загружаем конфиг из config.yaml c обработкой ошибок
	config, err := config.LoadConfig(pathConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем существует ли база, если нет, то создаём
	pathDb := config.DbFile
	var db *database.DataBase
	if _, err := os.Stat(pathDb); err != nil {
		fmt.Println(config.DbFile + " not exist")
		db = utilsDb.CreateDatabase(config.SourceUrl, config.DbFile, pathDb)
	} else {
		fmt.Println(config.DbFile + " exist")
		db = database.NewDataBase(pathDb)
	}

	// Обрабатываем флаги -o и -n
	if needShowDb {
		utilsDb.ShowDb(numComics, db)
	}
}
