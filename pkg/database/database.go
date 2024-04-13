package database

import (
	"encoding/json"
	"os"
)

// Представление даннх в бд
type ComicsDb struct {
	Image    string   `json:"url"`
	KeyWords []string `json:"keywords"`
}

// Бд для хранения комиксов
type DataBase struct {
	Path string
}

// Конструктор
func NewDataBase(path string) *DataBase {
	return &DataBase{
		Path: path,
	}
}

// Добавляет комикс в бд
func (DB *DataBase) AddOneComics(Id int, comics ComicsDb) {
	var mapDb map[int]ComicsDb

	// если добавляем первый комикс, значит файл пуст и надо создать мапу
	if Id == 1 {
		mapDb = make(map[int]ComicsDb)
	} else {
		mapDb = DB.GetDatabase()
	}

	// добавили комикс
	mapDb[Id] = comics

	// записали комикс
	bytes, _ := json.MarshalIndent(mapDb, "", "\t")
	os.WriteFile(DB.Path, bytes, 0644)
}

// Получаем бд из файла в виде мапы, нельзя вызывать если файла не существует
func (DB *DataBase) GetDatabase() map[int]ComicsDb {
	db := map[int]ComicsDb{}

	// читаем файл
	bytes, err := os.ReadFile(DB.Path)
	if err != nil {
		panic(err)
	}

	// заполняем мапу
	err = json.Unmarshal(bytes, &db)
	if err != nil {
		panic(err)
	}

	return db
}
