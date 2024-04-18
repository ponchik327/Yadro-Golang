package utilsDb

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/database"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/words"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/xkcd"
)

// Парсит флаги
func ParseFlags() string {
	pathConfig := flag.String("c", "config.yaml", "path to config file")
	flag.Parse()

	return *pathConfig
}

// Структура для удобства передачи ресурсов в воркер
type Resources struct {
	DataBase *database.DataBase
	Client   *xkcd.Client
}

// Задача, которую выполянет воркер выведена в отдельную функцию, которая ему передаётся
type funcWork func(id int, db *database.DataBase, client *xkcd.Client)

// В будущем передадим в воркер
func work(id int, db *database.DataBase, client *xkcd.Client) {
	// парсим с клиента
	comics := client.GetComicsById(id)
	// если не пустой
	if comics.Id != 0 {
		// конвертируем комикс спаршенный в представление бд
		comicsDb := convertComics(comics)
		// добавляем его в бд
		db.AddOneComics(id, comicsDb)
	}
}

// Создаёт бд
func CreateDatabase(sourceUrl string, pathDb string, numGorutine int) (*database.DataBase, error) {
	// создаём клиент и делаем запрос на id последнего комикса
	client := xkcd.NewClient(sourceUrl)
	idLastComics := client.GetIdLastComics()

	// открываем бд
	db, err := database.Open(pathDb)
	if err != nil {
		return nil, fmt.Errorf("error open database: %w", err)
	}

	ctx := context.Background()
	// добавялем контексту отслеживание сигналов
	// если придёт один из сигналов, он сделает вызов cancel()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	// штатный вызов cancel() в конце функции
	defer cancel()

	resources := Resources{
		DataBase: db,
		Client:   client,
	}

	// запускаем воркер-пул, передаём контекст, кол-во воркеров, максимальное кол-во задач, необходимые ресурсы и функцию с задачей
	workerPool(ctx, numGorutine, idLastComics, resources, work)

	fmt.Println(filepath.Base(pathDb) + " create")
	return db, nil
}

// Запускает паралельную обработку
func workerPool(ctx context.Context, countWorkers int, countTasks int, resources Resources, work funcWork) {
	wg := &sync.WaitGroup{}
	// по этому каналу будем подавать айди комиксов котрые надо обработать
	idComics := make(chan int, countTasks)
	// если обработка прошла успешно в этот канал придёт пустая структура
	result := make(chan struct{}, countTasks)

	// запускаем воркеры в отдельных горутнах, добавляю каждому задачу в ваит группе
	for i := 0; i < countWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, idComics, result, resources, work)
		}()
	}

	// эта функция через канал отправляет id комиксов, которые надо обработать
	giveTaks(countTasks, resources.DataBase, idComics)

	// эта горутина отвечает за вывод в консоль уведоление о скачки 100 комиксов
	go writeNotifications(ctx, result)

	wg.Wait()
	close(result)
}

// Даёт воркерам id, также проверяет есть ли данный ключ в бд, если есть, то не обрабатываем
func giveTaks(countTasks int, db *database.DataBase, idComics chan<- int) {
	for id := 1; id <= countTasks; id++ {
		_, ok := db.DataMap[id]
		if !ok {
			idComics <- id
		}
	}
	close(idComics)
}

// Каждые 100 обработанных комиксов оповещает пользователя
func writeNotifications(ctx context.Context, result <-chan struct{}) {
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-result:
			if !ok {
				return
			}

			count++
			if count%100 == 0 {
				fmt.Println("load " + strconv.Itoa(count) + " comics")
			}
		}
	}
}

// Выполняет работу, также при отмене контекста заканчивает выполнение горутины
func worker(ctx context.Context, idComics <-chan int, result chan<- struct{}, resources Resources, work funcWork) {
	for {
		select {
		// заканчиваем работу
		case <-ctx.Done():
			return
		// полезная работа
		case id, ok := <-idComics:
			// если прочитали из закрытого канала выходим, задач - нет
			if !ok {
				return
			}

			// работаем ...
			work(id, resources.DataBase, resources.Client)

			// пишем в канал, что обработали 1 комикс
			result <- struct{}{}
		}
	}
}

// Конвертируем из xkcd.Comics в database.ComicsDb
func convertComics(comics xkcd.Comics) database.ComicsDb {
	return database.ComicsDb{
		Image:    comics.Image,
		KeyWords: words.StemWords(strings.Fields(comics.Transcript + " " + comics.Alternative)),
	}
}
