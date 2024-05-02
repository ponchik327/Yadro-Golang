package bench_test

import (
	"log"
	"testing"

	"github.com/ponchik327/Yadro-Golang/tree/main/internal/config"
	"github.com/ponchik327/Yadro-Golang/tree/main/internal/search"
	"github.com/ponchik327/Yadro-Golang/tree/main/internal/utilsDb"
	"github.com/ponchik327/Yadro-Golang/tree/main/pkg/index"
)

func BenchmarkDefaultSearch(b *testing.B) {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("error open config: " + err.Error())
	}

	db, err := utilsDb.CreateDatabase(config.SourceUrl, config.DbFile, config.NumGorutine)
	if err != nil {
		log.Fatal("error create database: " + err.Error())
	}
	defer db.Close()

	index, err := index.CreateIndex(config.IndexFile, db)
	if err != nil {
		log.Fatal("error create index: " + err.Error())
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		search.Search("I'm following your questions", db, index, false, false)
		b.StopTimer()
	}
}

func BenchmarkIndexSearch(b *testing.B) {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("error open config: " + err.Error())
	}

	db, err := utilsDb.CreateDatabase(config.SourceUrl, config.DbFile, config.NumGorutine)
	if err != nil {
		log.Fatal("error create database: " + err.Error())
	}
	defer db.Close()

	index, err := index.CreateIndex(config.IndexFile, db)
	if err != nil {
		log.Fatal("error create index: " + err.Error())
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		search.Search("I'm following your questions", db, index, true, false)
		b.StopTimer()
	}
}
