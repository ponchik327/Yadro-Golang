package search

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/database"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/index"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/words"
)

func stemAndMappingRequest(request string) map[string]struct{} {
	wordsRequest := strings.Fields(request)

	mapRequset := make(map[string]struct{})
	for _, word := range words.StemWords(wordsRequest) {
		mapRequset[word] = struct{}{}
	}

	return mapRequset
}

func computeRelevanceDefault(mapRequset map[string]struct{}, db *database.DataBase) map[int]int {
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

func computeRelevanceIndex(mapRequset map[string]struct{}, index *index.IndexFile) map[int]int {
	idByRelevance := make(map[int]int)
	for word := range mapRequset {
		comicsIds := index.GetComicsIdByWord(word)
		for _, id := range comicsIds {
			idByRelevance[id]++
		}
	}

	return idByRelevance
}

func sortMapByRelevance(idByRelevance map[int]int) PairList {
	pl := make(PairList, len(idByRelevance))
	i := 0
	for k, v := range idByRelevance {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
}

type Pair struct {
	Id        int
	Relevance int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Relevance > p[j].Relevance }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func printCountRelevantComics(sortArrayByRelevance PairList, count int) {
	for i := 0; i < count; i++ {
		fmt.Print(i + 1)
		fmt.Print(": https://xkcd.com/")
		fmt.Print(sortArrayByRelevance[i].Id)
		fmt.Print("/info.0.json ")
		fmt.Print(sortArrayByRelevance[i].Relevance)
		fmt.Print("\n")
	}
}

func Search(request string, db *database.DataBase, index *index.IndexFile, isIndexSearch bool) {
	mapRequset := stemAndMappingRequest(request)

	var idByRelevance map[int]int
	if isIndexSearch {
		idByRelevance = computeRelevanceIndex(mapRequset, index)
	} else {
		idByRelevance = computeRelevanceDefault(mapRequset, db)
	}

	sortArrayByRelevance := sortMapByRelevance(idByRelevance)

	printCountRelevantComics(sortArrayByRelevance, 10)
}
