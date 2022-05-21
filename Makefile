run:
	go run .

build:
	go build -o app .
	./app

test:
	go test

clear:
	rm -R app