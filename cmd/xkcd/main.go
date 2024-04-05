package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/xkcd"
	"gopkg.in/yaml.v2"
)

type Config struct {
	SourceUrl string `yaml:"source_url"`
	DbFile    string `yaml:"db_file"`
}

func loadConfig(path string) (Config, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return Config{}, errors.New("unable to load config file")
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return Config{}, errors.New("yaml decode error")
	}

	return config, nil
}

func parseFlags() (bool, int) {
	needShowDb := flag.Bool("o", false, "display db")
	numComics := flag.Int("n", -1, "count comics to display")

	flag.Parse()

	return *needShowDb, *numComics
}

func main() {
	rootDir := filepath.Join("..", "..")

	pathConfig := filepath.Join(rootDir, "config.yaml")
	config, err := loadConfig(pathConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)

	needShowDb, numComics := parseFlags()

	if needShowDb {
		if numComics == -1 {
			fmt.Println(numComics)
		} else {
			for i := 0; i < numComics; i++ {
				fmt.Println(i)
			}
		}
	}

	pathDb := filepath.Join(rootDir, config.DbFile)
	if _, err := os.Stat(pathDb); err != nil {
		fmt.Println("file not exist")
		allComics, err := xkcd.ParseComics(config.SourceUrl)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(allComics)
	} else {
		fmt.Println("file exist")
	}
}
