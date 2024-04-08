package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
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
	fmt.Println(path)
	bytes, err := os.ReadFile(path)

	if err != nil {
		return Config{}, fmt.Errorf("unable to load config file, err: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return Config{}, fmt.Errorf("yaml decode error, err: %w", err)
	}

	return config, nil
}

// Парсит флаги
func parseFlags() (string, bool, int) {
	pathConfig := flag.String("c", "config.yaml", "path to config file")
	needShowDb := flag.Bool("o", false, "display db")
	numComics := flag.Int("n", math.MaxInt, "count comics to display")

	flag.Parse()

	return *pathConfig, *needShowDb, *numComics
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
	if numComics == math.MaxInt {
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
	pathConfig, needShowDb, numComics := parseFlags()

	// Загружаем конфиг из config.yaml c обработкой ошибок
	config, err := loadConfig(pathConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем существует ли база, если нет, то создаём
	pathDb := config.DbFile
	if _, err := os.Stat(pathDb); err != nil {
		fmt.Println(config.DbFile + " not exist")
		createDatabase(config, pathDb)
	} else {
		fmt.Println(config.DbFile + " exist")
	}

	// Обрабатываем флаги -o и -n
	if needShowDb {
		showDb(numComics, pathDb)
	}
}
