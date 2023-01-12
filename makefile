build:
	go build -o server main.go

run: build
	./server

watch:
	ulimit -n 1000 #increase the file watch limit, might required on MacOS
	go run ~/go/pkg/mod/github.com/cespare/reflex@v0.3.1/ -s -r '\.go$$' make run