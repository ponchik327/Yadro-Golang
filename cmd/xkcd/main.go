package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/database"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/xkcd"
	"gopkg.in/yaml.v2"
)

type Config struct {
	SourceUrl string `yaml:"source_url"`
	DbFile    string `yaml:"db_file"`
}

// Читает конфиг
func loadConfig(path string) (Config, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return Config{}, errors.New("unable to load config file")
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return Config{}, errors.New("yaml decode error")
	}

	return config, nil
}

const Default int = 0

// Парсит флаги
func parseFlags() (bool, int) {
	needShowDb := flag.Bool("o", false, "display db")
	numComics := flag.Int("n", Default, "count comics to display")

	flag.Parse()

	return *needShowDb, *numComics
}

// Печатает комикс
func printComics(num int, comics *database.ComicsDb, id int) {
	fmt.Println(strconv.Itoa(num) + " comics")
	fmt.Println("id: " + strconv.Itoa(id))
	fmt.Println("image: " + comics.Image)
	fmt.Println("keywords: ")
	fmt.Println(comics.KeyWords)
	fmt.Println("--------------------------------------------")
}

// Создаёт бд
func createDatabase(config Config, pathDb string) {
	// парсим комиксы в массив
	allComics, err := xkcd.ParseComics(config.SourceUrl)
	if err != nil {
		log.Fatal(err)
	}

	// передаёи массив для создания database.json
	database.CreateDatabase(allComics, pathDb)
	fmt.Println(config.DbFile + " create")
}

// Отображаем записи из бд в нужном количестве
func showDb(numComics int, pathDb string) {
	// достаём всю базу данных в виде мапы
	db, err := database.GetDatabase(pathDb)
	if err != nil {
		log.Fatal(err)
	}

	// в зависимости от флага выводим нужное количество комиксов
	if numComics == Default {
		i := 1
		// выводим всё комиксы
		for id, comics := range db {
			printComics(i, &comics, id)
			i++
		}
	} else {
		// проверям, что количество которое надо вывести меньше, чем размер мапы в бд
		countComics := math.Min(float64(numComics), float64(len(db)))
		i := 1
		for id, comics := range db {
			// выводим countComics комиксов
			if i <= int(countComics) {
				printComics(i, &comics, id)
			}
			i++
		}
	}
}

func main() {
	rootDir := filepath.Join("..", "..")

	// Загружаем конфиг из config.yaml c обработкой ошибок
	pathConfig := filepath.Join(rootDir, "config.yaml")
	config, err := loadConfig(pathConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем существует ли база, если нет, то создаём
	pathDb := filepath.Join(rootDir, config.DbFile)
	if _, err := os.Stat(pathDb); err != nil {
		fmt.Println(config.DbFile + " not exist")
		createDatabase(config, pathDb)
	} else {
		fmt.Println(config.DbFile + " exist")
	}

	// Считываем флаги и обрабатываем их
	needShowDb, numComics := parseFlags()
	if needShowDb {
		showDb(numComics, pathDb)
	}
}
