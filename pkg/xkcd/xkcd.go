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
func (c *Client) GetIdLastComics() int {
	urlLastComics, err := url.JoinPath(c.SourceUrl, c.EndPoints["GetComics"])
	if err != nil {
		panic(err)
	}

	lastComics := c.parseOneComics(urlLastComics)

	return lastComics.Id
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
