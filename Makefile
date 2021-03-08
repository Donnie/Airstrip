builddev:
	docker-compose build --pull

dev:
	docker-compose --env-file ./.env.local up

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

# CI
build:
	@echo "Building for prod"
	docker build -t donnieashok/airstrip:prod .

deploy: build
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	docker push donnieashok/airstrip:prod
	@echo "Deployed!"

# Prod
live:
	ssh root@airstrip docker pull donnieashok/airstrip:prod
	- ssh root@airstrip docker stop airstrip
	- ssh root@airstrip docker rm airstrip
	scp -r ./.env root@airstrip:/root/
	ssh root@airstrip docker run -d --restart on-failure --env-file /root/.env --name airstrip donnieashok/airstrip:prod
	ssh root@airstrip rm /root/.env
	@echo "Is live"
