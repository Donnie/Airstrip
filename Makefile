build:
	docker-compose build --pull

builddev:
	docker-compose -f dev-compose.yml build --pull

dev:
	docker-compose -f dev-compose.yml up

up:
	docker-compose up

sql:
	docker-compose run -e PGPASSWORD=postgres postgres psql --host=postgres --username=airstrip --dbname=airstrip

clean:
	@echo "Cleaning Docker environment..."
	docker-compose stop
	docker-compose down -v
