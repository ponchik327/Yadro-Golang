build: 
	go build -o  cmd/xkcd/xkcd.exe cmd/xkcd/main.go

run_exe: 
	cmd/xkcd/xkcd.exe -o -n 2

run: 
	go run cmd/xkcd/main.go

run_with_arg: 
	go run cmd/xkcd/main.go -o -n 3