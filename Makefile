default: build docker docker-up

build:
	go build
	GOOS=linux go build -o app

docker:
	docker-compose build app

docker-up:
	docker-compose up app
