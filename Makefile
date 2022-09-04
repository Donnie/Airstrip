builddev:
	docker build -f Dockerfile.dev -t donnieashok/airstrip:dev .

dev:
	docker-compose --env-file ./.env.local up

sql:
	sqlite lite/sql.db

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
	scp ./.env root@airstrip:/home/airstrip/
	scp ./db/sql.db root@airstrip:/home/airstrip/db/
	ssh root@airstrip docker run -d --restart on-failure -v /home/airstrip/db:/db --env-file /home/airstrip/.env --name airstrip donnieashok/airstrip:prod
	ssh root@airstrip rm /home/airstrip/.env
	@echo "Is live"
