package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Comics struct {
	Id          int    `json:"num"`
	Image       string `json:"img"`
	Transcript  string `json:"transcript"`
	Alternative string `json:"alt"`
}

func parseOneComics(endpoint string) (Comics, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return Comics{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Comics{}, err
	}

	var comics Comics
	err = json.Unmarshal(body, &comics)
	if err != nil {
		return Comics{}, err
	}

	return comics, nil
}

func ParseComics(sourceUrl string) ([]Comics, error) {

	lastComics, err := url.JoinPath(sourceUrl, "/info.0.json")
	if err != nil {
		return []Comics{}, err
	}

	fmt.Println(lastComics)

	comics, err := parseOneComics(lastComics)
	if err != nil {
		return []Comics{}, err
	}

	var allComics []Comics
	allComics = append(allComics, comics)

	return allComics, nil
}
