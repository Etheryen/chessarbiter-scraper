FLAGS = 

dev:
	clear && gow -s -v run cmd/main.go $(FLAGS)

build:
	go build -o bin/app cmd/main.go

run: build
	./bin/app $(FLAGS)

test:
	go test -v ./... -count=1
