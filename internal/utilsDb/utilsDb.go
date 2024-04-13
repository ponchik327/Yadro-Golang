package utilsDb

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/database"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/words"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/xkcd"
)

// Парсит флаги
func ParseFlags() (string, bool, int) {
	pathConfig := flag.String("c", "config.yaml", "path to config file")
	needShowDb := flag.Bool("o", false, "display db")
	numComics := flag.Int("n", math.MaxInt, "count comics to display")

	flag.Parse()

	return *pathConfig, *needShowDb, *numComics
}

// Создаёт бд
func CreateDatabase(sourceUrl string, dbFile string, pathDb string) *database.DataBase {
	// создаём клиент и делаем запрос на id последнего комикса
	client := xkcd.NewClient(sourceUrl)
	idLastComics := client.GetIdLastComics()

	// в цикле проходимся по всем id и заполняем бд
	db := database.NewDataBase(pathDb)
	for id := 1; id <= idLastComics; id++ {
		// каждые 100 комиксов пишем прогресс
		if id%100 == 0 {
			fmt.Println("load " + strconv.Itoa(id) + " comics")
		}

		// парсим с клиента
		comics := client.GetComicsById(id)
		// если не пустой
		if comics.Id != 0 {
			// конвертируем комикс спаршенный в представление бд
			comicsDb := convertComics(comics)
			// добавляем его в бд
			db.AddOneComics(id, comicsDb)
		}
	}

	fmt.Println(dbFile + " create")
	return db
}

// Отображаем записи из бд в нужном количестве
func ShowDb(numComics int, dataBase *database.DataBase) {
	// достаём всю базу данных в виде мапы
	mapDb := dataBase.GetDatabase()

	// в зависимости от флага выводим нужное количество комиксов
	if numComics == math.MaxInt {
		i := 1
		// выводим всё комиксы
		for id, comics := range mapDb {
			printComics(i, &comics, id)
			i++
		}
	} else {
		// проверям, что количество которое надо вывести меньше, чем размер мапы в бд
		countComics := math.Min(float64(numComics), float64(len(mapDb)))
		i := 1
		for id, comics := range mapDb {
			// выводим countComics комиксов
			if i <= int(countComics) {
				printComics(i, &comics, id)
			}
			i++
		}
	}
}

// Форматированный вывод для комикса
func printComics(num int, comics *database.ComicsDb, id int) {
	fmt.Println(strconv.Itoa(num) + " comics")
	fmt.Println("id: " + strconv.Itoa(id))
	fmt.Println("image: " + comics.Image)
	fmt.Println("keywords: ")
	fmt.Println(comics.KeyWords)
	fmt.Println("--------------------------------------------")
}

// Конвертируем из xkcd.Comics в database.ComicsDb
func convertComics(comics xkcd.Comics) database.ComicsDb {
	return database.ComicsDb{
		Image:    comics.Image,
		KeyWords: words.StemWords(strings.Fields(comics.Transcript + " " + comics.Alternative)),
	}
}
