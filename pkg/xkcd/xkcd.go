package xkcd

import (
	"encoding/json"
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

// Клиент для парса комиксов
type Client struct {
	SourceUrl string
	EndPoints map[string]string
}

// Контсруктор
func NewClient(sourceUrl string) *Client {
	// список возможных путей
	endPoints := map[string]string{
		"GetComics": "info.0.json",
	}

	return &Client{
		SourceUrl: sourceUrl,
		EndPoints: endPoints,
	}
}

// Возвращает комикс если такой можно спарсить, в пративном случае пустой комикс
func (c *Client) GetComicsById(IdComics int) Comics {
	endpoint := c.makeEndpoint(IdComics)
	comics := c.parseOneComics(endpoint)

	return comics
}

// Id последнего комикса
// Теоретически работает за O(log(n))
func (c *Client) GetIdLastComics() int {

	// c помощью экспоненциального поиска ищем промежуток в котором находится id последнего комикса
	id := 1
	idPrev := id
	for c.GetComicsById(id).Id != 0 && id != 404 {
		idPrev = id
		id *= 2
	}

	// с помощью бинарного поиска ищем на этом отрезке нужный id
	id = lowerBound(c, idPrev, id) - 1

	return id
}

// Реализация бинарного поиска для нашего случая
// Можем представить, что у нас есть отсортированный массив [0 0 0 1 1 1], где 0 - комикс существует, а 1 - комикс не существует
// И наш алгоритм ищет перое вхождение 1, то есть id, который следует за id последнего комикса
func lowerBound(c *Client, low int, high int) int {
	mid := 0
	for low <= high {
		mid = (low + high) / 2

		var num int
		if c.GetComicsById(mid).Id == 0 {
			num = 1
		} else {
			num = 0
		}

		if num >= 1 {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return low
}

// Парсит один комикс
func (c *Client) parseOneComics(endpoint string) Comics {
	// делаем get-запрос
	resp, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// если статус ответа 200, то заполняем Comics
	var comics Comics
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&comics)
		if err != nil {
			panic(err)
		}
	}

	return comics
}

// Подготавливает url для получения комиксов по их id
func (c *Client) makeEndpoint(IdComics int) string {
	urlComics, err := url.JoinPath(c.SourceUrl, strconv.Itoa(IdComics), c.EndPoints["GetComics"])
	if err != nil {
		panic(err)
	}

	return urlComics
}
