psql:
	docker-compose exec postgres psql -U db_user workshop

test:
	go test ./...

format:
	go fmt ./...

generate:
	go generate ./...

up:
	docker-compose up -d

stop:
	docker-compose stop

down:
	docker-compose down

logs:
	docker-compose logs app
