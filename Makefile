server:
	go build api/main.go

build:
	go build -o bin/server api/main.go

d.up:
	docker-compose up

d.down:
	docker-compose down

d.up.build:
	docker-compose up --build
