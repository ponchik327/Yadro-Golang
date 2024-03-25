build: 
	go build -o bin/myapp myapp.go

run: 
	go run myapp.go

run_with_arg_1:
	go run myapp.go -s "follower brings bunch of questions"

run_with_arg_2:
	go run myapp.go -s "i'll follow you as long as you are following me"