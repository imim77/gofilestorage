build:
	go build -o bin/gofilestorage
run: build
	./bin/gofilestorage
test:
	go test ./... -v 
