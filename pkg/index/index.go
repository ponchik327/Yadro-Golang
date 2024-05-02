package index

import (
	"encoding/json"
	"os"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/database"
)

// Индексный файл
type IndexFile struct {
	Path           string
	WordToComicsId map[string][]int
}

// Использую бд создает индексный файл и записывает его на диск
func CreateIndex(path string, db *database.DataBase) (*IndexFile, error) {
	index := IndexFile{
		Path:           path,
		WordToComicsId: make(map[string][]int),
	}

	for id, comics := range db.DataMap {
		for _, keyWord := range comics.KeyWords {
			index.AddIndexByWord(keyWord, id)
		}
	}

	bytes, err := json.MarshalIndent(index.WordToComicsId, "", "\t")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(index.Path, bytes, 0644)
	if err != nil {
		return nil, err
	}

	return &index, nil
}

// Добавляет в индекс новый id по слову
func (index *IndexFile) AddIndexByWord(word string, id int) {
	comicsIds := index.WordToComicsId[word]
	comicsIds = append(comicsIds, id)
	index.WordToComicsId[word] = comicsIds
}

// Выдаёт все комиксы из индекса связанные с ключевым словом
func (index *IndexFile) GetComicsIdByWord(word string) []int {
	return index.WordToComicsId[word]
}
