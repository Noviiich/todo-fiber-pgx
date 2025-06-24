APP_NAME = todoapp
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/platform/migrations
DATABASE_URL = postgres://postgres:password@todoapp-postgres/postgres?sslmode=disable

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

docker.run: docker.network docker.postgres swag docker.fiber migrate.up

docker.network:
	docker network inspect todoapp-network >/dev/null 2>&1 || \
	docker network create -d bridge todoapp-network

docker.fiber.build:
	docker build -t fiber .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name todoapp-fiber \
		--network todoapp-network \
		-p 5000:5000 \
		fiber

docker.postgres:
	docker run --rm -d \
		--name todoapp-postgres \
		--network todoapp-network \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=postgres \
		-v ${HOME}/todoapp-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.stop: docker.stop.fiber docker.stop.postgres

docker.stop.fiber:
	docker stop todoapp-fiber

docker.stop.postgres:
	docker stop todoapp-postgres

swag:
	swag init