.PHONY: build up down restart clear test migrate-up migrate-down

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose down	
	docker-compose up -d

clear:
	docker-compose down -v --remove-orphans
	docker container prune -f
	docker image prune -f
	docker volume prune -f
	@echo "======================="	

test:
	go test -v ./...