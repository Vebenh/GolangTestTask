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

run:
	mkdir -p ../bin/config
	go build -o ../bin/main ./cmd/main.go
	cp ./config/*.yaml ../bin/config
	cd ../bin && ./main