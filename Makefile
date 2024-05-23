
graph.gen:
	go run github.com/99designs/gqlgen generate --verbose

docker.rebuild:
	docker compose --env-file ./docker.env up -d --build app

docker.run:
	docker compose --env-file ./docker.env up -d

docker.run.db:
	docker compose --env-file ./docker.env up -d postgres

docker.run.migrate:
	docker compose --env-file ./docker.env up -d migrate

docker.down:
	docker compose down

migrate.up:
	migrate -path ./migrations -database "postgres://yks:yksadm@localhost:5432/postgres?sslmode=disable" up

migrate.down:
	migrate -path ./migrations -database "postgres://yks:yksadm@localhost:5432/postgres?sslmode=disable" down

tests.run:
	go test ./...

tests.cover:
	go test -cover ./...