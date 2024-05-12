build: 
	go build -o cmd/xkcd/xkcd.exe cmd/xkcd/main.go

run_exe: 
	cmd/xkcd/xkcd.exe 

run: 
	go run cmd/xkcd/main.go

run_search: 
	go run cmd/xkcd/main.go -s "I'm following your questions"

run_index_search:
	go run cmd/xkcd/main.go -s "I'm following your questions" -i

run_bencmarks:
	go test -bench=.