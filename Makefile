.PHONY: build up down clear test migrate-up migrate-down

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down -v --remove-orphans

clear:
	docker-compose down -v --remove-orphans
	docker container prune -f
	docker image prune -f
	docker volume prune -f
	@echo "======================="	

test:
	go test -v ./...