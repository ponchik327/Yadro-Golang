package database

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/words"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/xkcd"
)

// Представление даннх в бд
type ComicsDb struct {
	Image    string   `json:"url"`
	KeyWords []string `json:"keywords"`
}

// Из вектора комиксов в представлении при парсинге (Comics)
// Длает мапу по id представлений комикса в бд (ComicsDb)
func CreateDatabase(comics []xkcd.Comics, pathDb string) {
	db := map[int]ComicsDb{}

	// пробразуем массив в мапу с нужными полями
	// попутно делаем стемминг транскрипциии и альтернативного представления
	for _, oneComics := range comics {
		db[oneComics.Id] = ComicsDb{
			Image:    oneComics.Image,
			KeyWords: words.StemWords(strings.Fields(oneComics.Transcript + " " + oneComics.Alternative)),
		}
	}

	// пишем получившуюся мапу в файл database.json
	bytes, _ := json.MarshalIndent(db, "", "\t")
	os.WriteFile(pathDb, bytes, 0644)
}

// Вытаскиваем из бд представление комиксов в виде мапы
func GetDatabase(pathDb string) (map[int]ComicsDb, error) {
	db := map[int]ComicsDb{}

	bytes, err := os.ReadFile(pathDb)
	if err != nil {
		return db, err
	}

	err = json.Unmarshal(bytes, &db)
	if err != nil {
		return db, err
	}

	return db, nil
}
