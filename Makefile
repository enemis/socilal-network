ifneq (,$(wildcard ./.env))
    include .env
    export
endif

IMAGE=social_network


all: run migrations-up
run:
	docker compose up -d
stop:
	docker compose down	

migrations-up:
	docker exec -i -w /usr/src/app/migrations social_network sh -c "goose postgres 'user=${DB_USERNAME} host=${DB_HOST} dbname=${DB_USERNAME} sslmode=${DB_SSLMODE} password=${DB_PASSWORD}' up"

migrations-down:
	docker exec -i -w /usr/src/app/migrations social_network sh -c "goose postgres 'user=${DB_USERNAME} host=${DB_HOST} dbname=${DB_USERNAME} sslmode=${DB_SSLMODE} password=${DB_PASSWORD}' down"


seed:
	docker exec -i social_network sh -c "/usr/local/bin/seeder"