package xkcd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Представление даннх при парсинге
type Comics struct {
	Id          int    `json:"num"`
	Image       string `json:"img"`
	Transcript  string `json:"transcript"`
	Alternative string `json:"alt"`
}

// Обработка одного комикса
func parseOneComics(endpoint string) (Comics, error) {
	// делаем запрос
	resp, err := http.Get(endpoint)
	if err != nil {
		return Comics{}, err
	}
	defer resp.Body.Close()

	// проверяем, что ответ пришёл корректно
	if resp.Status != "200 OK" {
		return Comics{}, err
	}

	// читаем ответ в массив байтов
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Comics{}, err
	}

	// преобразуем ответ в структуру
	var comics Comics
	err = json.Unmarshal(body, &comics)
	if err != nil {
		return Comics{}, err
	}

	return comics, nil
}

const endPointGetComics = "info.0.json"

// Парсим последний комикс отдельно, чтобы узнать сколько их всего
func parseLastComics(sourceUrl string) (Comics, error) {
	urlLastComics, err := url.JoinPath(sourceUrl, endPointGetComics)
	if err != nil {
		return Comics{}, err
	}

	lastComics, err := parseOneComics(urlLastComics)
	if err != nil {
		return Comics{}, err
	}

	return lastComics, nil
}

// Парсим все комиксы со страницы
func ParseComics(sourceUrl string) ([]Comics, error) {

	// парсим послендий для понимания, сколько их всего
	lastComics, err := parseLastComics(sourceUrl)
	if err != nil {
		return []Comics{}, err
	}

	var allComics []Comics
	allComics = append(allComics, lastComics)

	// добавляем все спаршенные комиксы в массив
	for i := 1; i < lastComics.Id; i++ {
		// задаём url
		urlComics, err := url.JoinPath(sourceUrl, strconv.Itoa(i), endPointGetComics)
		if err != nil {
			return []Comics{}, err
		}

		// парсим одну штуку
		comics, err := parseOneComics(urlComics)
		if err != nil {
			return []Comics{}, err
		}

		// добавляем в массив
		if comics.Id != 0 {
			allComics = append(allComics, comics)
		}
	}

	return allComics, nil
}
