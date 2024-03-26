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

// Проходит по строке и удаляет символы, которые не являются буквами
func deleteNonLetters(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func main() {
	// Прасим аргументы из командной строки
	wordsPtr := flag.String("s", "a dog", "a string to normalization")
	flag.Parse()

	// Открываем файл со стоп-словами https://countwordsfree.com/stopwords?ysclid=lu7cl0fbo0793135203
	// В нём уже содеражться популярные слова с апострофами
	file, err := os.Open("stop_words_english.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Читаем стоп-слова из файла
	// и добавляем их в map без ассоциации (используем map как set)
	var stopWords = make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stopWords[scanner.Text()] = true
	}

	// Будем использовать для того, чтобы хранить встретившиеся слова
	var encountWords = make(map[string]bool)

	// Перебираем строки(слова) переданные как аргументы,
	// если они не являются стоп-словами и не встречались ранее,
	// то нормализуем их и выводим в консоль
	for _, word := range strings.Fields(*wordsPtr) {
		_, isStopWord := stopWords[word]
		if !isStopWord {
			stemmed, err := snowball.Stem(deleteNonLetters(word), "english", true)
			_, isEncount := encountWords[stemmed]
			if (err == nil) &&
				(!isEncount) {
				fmt.Print(stemmed, " ")
			}
			encountWords[stemmed] = true
		}
	}
}
