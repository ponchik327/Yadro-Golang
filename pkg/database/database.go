package database

import (
	"encoding/json"
	"os"
	"sync"
)

// Представление даннх в бд
type ComicsDb struct {
	Image    string   `json:"url"`
	KeyWords []string `json:"keywords"`
}

// Бд для хранения комиксов
type DataBase struct {
	Path    string
	Mutex   sync.Mutex
	DataMap map[int]ComicsDb
}

// Открывает бд,
// если файл существует, то считывает его в мапу структуры,
// если нет, то создаёт пустой
func Open(path string) (*DataBase, error) {
	// создали пустую структуру
	db := DataBase{
		Path: path,
	}

	// проверили существование файла
	if _, err := os.Stat(db.Path); err == nil {
		// считали файл
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		// переложили данные в мапу структуры
		err = json.Unmarshal(bytes, &db.DataMap)
		if err != nil {
			return nil, err
		}
	} else {
		// инициализировали пустую мапу
		db.DataMap = make(map[int]ComicsDb)

		// преобразовали в файл
		bytes, err := json.MarshalIndent(db.DataMap, "", "\t")
		if err != nil {
			return nil, err
		}

		// созадли пустой файл бд
		err = os.WriteFile(db.Path, bytes, 0644)
		if err != nil {
			return nil, err
		}
	}

	return &db, nil
}

// Закрывает бд, сохраняет мапу стуктуры в файл
func (DB *DataBase) Close() error {
	bytes, err := json.MarshalIndent(DB.DataMap, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(DB.Path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Добавляет комикс в бд
func (DB *DataBase) AddOneComics(id int, comics ComicsDb) {
	DB.Mutex.Lock()
	defer DB.Mutex.Unlock()
	DB.DataMap[id] = comics
}
