package words

import (
	"bufio"
	"os"
	"strings"

	"github.com/kljensen/snowball"
)

// Загрузка стоп-слов из файла stop_words_english.txt в мапу, для удобства проверки
func loadStopWords(path string) map[string]bool {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var stopWords = make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stopWords[scanner.Text()] = true
	}

	return stopWords
}

// Вспомогательная функция для чистки слов
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

// Реализация стемминга из прошлого задания, вынес в функцию для просты использования
func StemWords(words []string) []string {

	stopWords := loadStopWords("stop_words_english.txt")

	// мапа для проверки повторяющихся слов
	var encountWords = make(map[string]bool)
	// массив очищенных слов
	stemmedWords := make([]string, 0)

	for _, word := range words {

		word = strings.ToLower(word)

		// первая проверка на стоп-слово
		_, isStopWord := stopWords[word]
		if !isStopWord {

			stemmed, err := snowball.Stem(deleteNonLetters(word), "english", true)

			// вторая проверка на стоп-слово
			// почему-то некоторые слова после стемминга начинали подходить под стоп-слова
			_, isStopWord = stopWords[stemmed]

			_, isEncount := encountWords[stemmed]

			if (err == nil) &&
				(!isEncount) &&
				(stemmed != "") &&
				(!isStopWord) {
				stemmedWords = append(stemmedWords, stemmed)
			}

			encountWords[stemmed] = true
		}
	}
	return stemmedWords
}
