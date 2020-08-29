build:
	docker-compose build --pull

up:
	docker-compose up

clean:
	@echo "Cleaning Docker environment..."
	docker-compose stop
	docker-compose down -v
	docker-compose rm airstrip
	docker-compose rm postgres
