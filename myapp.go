package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	// Библиотека для стемминга(нормализации) слов
	"github.com/kljensen/snowball"
)

func main() {
	// Прасим аргументы из командной строки
	wordsPtr := flag.String("s", "a dog", "a string to normalization")
	flag.Parse()

	// Открываем файл со стоп-словами https://countwordsfree.com/stopwords?ysclid=lu7cl0fbo0793135203
	file, err := os.Open("stop_words_english.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Читаем стоп-слова из файла
	// и добавляем их в map без ассоциации (используем map как set)
	var stop_words = make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stop_words[scanner.Text()] = struct{}{}
	}

	// Перебираем строки переданные как аргументы,
	// если они не являются стоп-словами,
	// то нормализуем их и выводим в консоль
	for _, word := range strings.Fields(*wordsPtr) {
		_, is_stop_word := stop_words[word]
		if !is_stop_word {
			stemmed, err := snowball.Stem(word, "english", true)
			if err == nil {
				fmt.Print(stemmed, " ")
			}
		}
	}
}
