#DC = docker-compose
#
#up:
#	$(DC) up --build -d
#
#down:
#	$(DC) down
#
#logs:
#	$(DC) logs -f app
#
#exec-app:
#	$(DC) exec app sh
#
#clean:
#	$(DC) down -v
#
#.PHONY: up down logs exec-app clean


# for local run
run:
	mkdir -p ../bin/config
	go build -o ../bin/main ./cmd/main.go
	cp ./config/*.yaml ../bin/config
	cd ../bin && ./main

build:
	docker-compose up --build

rebuild:
	docker-compose down
	docker-compose up --build -d
