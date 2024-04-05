package words

import (
	"bufio"
	"os"
	"strings"

	"github.com/kljensen/snowball"
)

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

func StemWords(words []string) []string {

	stopWords := loadStopWords("stop_words_english.txt")

	var encountWords = make(map[string]bool)
	stemmedWords := make([]string, 0)

	for _, word := range words {

		word = strings.ToLower(word)

		_, isStopWord := stopWords[word]
		if !isStopWord {

			stemmed, err := snowball.Stem(deleteNonLetters(word), "english", true)

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
