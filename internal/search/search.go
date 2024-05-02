package search

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/database"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/index"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/words"
)

// Преобразует поисковый запрос в мапу отстемленных слов
func stemAndMappingRequest(request string) map[string]struct{} {
	wordsRequest := strings.Fields(request)

	mapRequset := make(map[string]struct{})
	for _, word := range words.StemWords(wordsRequest) {
		mapRequset[word] = struct{}{}
	}

	return mapRequset
}

// Вычисление релевантности запроса с помошью обращения к бд
func computeRelevanceDefault(mapRequset map[string]struct{}, db *database.DataBase) map[int]int {
	idByRelevance := make(map[int]int)

	// для каждого коимкса из бд перебираем все ключевые слова
	for id, comics := range db.DataMap {
		for _, keyWord := range comics.KeyWords {
			// если это слово есть в посиковом запросе увеличиваем релевантность
			_, ok := mapRequset[keyWord]
			if ok {
				idByRelevance[id]++
			}
		}
	}

	return idByRelevance
}

// Вычисление релевантности запроса с помошью подготовленного индексного файла
func computeRelevanceIndex(mapRequset map[string]struct{}, index *index.IndexFile) map[int]int {
	idByRelevance := make(map[int]int)

	// проходимся по словам посикового запроса
	for word := range mapRequset {
		// из индекса вытаскиваем комиксы которые содержат это слово
		comicsIds := index.GetComicsIdByWord(word)
		// у этих комиксов увеличиваем релеватность
		for _, id := range comicsIds {
			idByRelevance[id]++
		}
	}

	return idByRelevance
}

// преобразуем мапу в массив и сортируем его в порядке убыванию релеватности
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

// --------------------- необходимо для сортировки массива ------------------------
type Pair struct {
	Id        int
	Relevance int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Relevance > p[j].Relevance }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ---------------------------------------------------------------------------------

// Печатает нужное количество комиксов в нужном формате
func printCountRelevantComics(sortArrayByRelevance PairList, count int, is_need_print bool) {
	for i := 0; i < count; i++ {
		// аргумент is_need_print нужен, чтобы во время запуска бенчмарка не выводились комиксы
		if is_need_print {
			fmt.Print(i + 1)
			fmt.Print(": https://xkcd.com/")
			fmt.Print(sortArrayByRelevance[i].Id)
			fmt.Print("/info.0.json ")
			fmt.Print(sortArrayByRelevance[i].Relevance)
			fmt.Print("\n")
		}
	}
}

// Осуществляет посик, выполняя все шаги по очереди
func Search(request string, db *database.DataBase, index *index.IndexFile, isIndexSearch bool, is_need_print bool) {
	mapRequset := stemAndMappingRequest(request)

	var idByRelevance map[int]int
	// реализует разные виды посика в зависимости от аргумента isIndexSearch
	if isIndexSearch {
		idByRelevance = computeRelevanceIndex(mapRequset, index)
	} else {
		idByRelevance = computeRelevanceDefault(mapRequset, db)
	}

	sortArrayByRelevance := sortMapByRelevance(idByRelevance)

	printCountRelevantComics(sortArrayByRelevance, 10, is_need_print)
}
