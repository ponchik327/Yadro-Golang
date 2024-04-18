package search

import (
	"strings"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/database"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/words"
)

func stemAndMappingRequest(request string) map[string]struct{} {
	wordsRequest := strings.Fields(request)

	mapRequset := make(map[string]struct{}, len(wordsRequest))
	for _, word := range words.StemWords(wordsRequest) {
		mapRequset[word] = struct{}{}
	}

	return mapRequset
}

func computeRelevance(mapRequset map[string]struct{}, db *database.DataBase) map[int]int {
	idByRelevance := make(map[int]int)
	for id, comics := range db.DataMap {
		for _, keyWord := range comics.KeyWords {
			_, ok := mapRequset[keyWord]
			if ok {
				idByRelevance[id]++
			}
		}
	}

	return idByRelevance
}

func sortMapByValue(wordFrequencies map[string]int) PairList{
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key string
	Value int
}

func printCountRelevantComics(idByRelevance map[int]int, count int) {
	for i := 0; i < count; i++ {

	}
}

func Search(request string, db *database.DataBase) {
	mapRequset := stemAndMappingRequest(request)

	idByRelevance := computeRelevance(mapRequset, db)

	printCountRelevantComics(idByRelevance, 10)
}
