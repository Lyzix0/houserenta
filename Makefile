include .env
export

export PROJECT_ROOT=$(CURDIR)

env-up:
	@docker compose up -d renta-app-postgres

env-down:
	@docker compose down renta-app-postgres

env-cleanup:
	@docker compose down renta-app-postgres; \
	sudo rm -rf ${PROJECT_ROOT}/out/pgdata; \
	echo "Files were cleaned!";

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Add param seq. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \

	@docker-compose run --rm renta-app-postgres-migrate \
		create \
		-ext sql \
		-dir //migrations \
		-seq "$(seq)"


migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Add param action. Example: make migrate-action action=down 2"; \
		exit 1; \
	fi
	@MSYS_NO_PATHCONV=1 docker compose run --rm renta-app-postgres-migrate \
		-path //migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@renta-app-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/app/main.go

ps:
	@docker compose ps

logs-cleanup:
	@sudo rm -rf ${PROJECT_ROOT}/out/logs; \
	echo "Logs were cleaned!";

app-deploy:
	@docker compose up -d --build renta-app

swagger-gen:
	@swag init -g cmd/app/main.go -o docs
	
swagger-build:
	@docker compose build --no-cache swagger