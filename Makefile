build:
	@echo "Building for prod"
	docker-compose build --pull

builddev:
	docker-compose -f dev-compose.yml build --pull

dev:
	docker-compose -f dev-compose.yml up

deploy: build
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	docker push donnieashok/airstrip:prod
	@echo "Deployed!"

up:
	@echo "Running for Prod"
	docker-compose up

sql:
	docker-compose run -e PGPASSWORD=postgres postgres psql --host=airstrip_db --username=airstrip --dbname=airstrip

dump:
	docker exec -e PGPASSWORD=postgres airstrip_db pg_dump --username=airstrip airstrip > airstrip.sql

migrate:
	docker exec -e PGPASSWORD=postgres -i airstrip_db psql --username airstrip --dbname airstrip < ./airstrip.sql

clean:
	@echo "Cleaning Docker environment..."
	docker-compose stop
	docker-compose down -v
