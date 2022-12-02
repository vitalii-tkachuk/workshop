psql:
	docker-compose exec postgres psql -U db_user workshop

table:
	docker-compose exec postgres psql -U db_user workshop sh -c 'CREATE TABLE IF NOT EXISTS users(id UUID, name VARCHAR)'

test:
	go test ./...

format:
	go fmt ./...

generate:
	go generate ./...
