builddev:
	docker build -f Dockerfile.dev -t donnieashok/airstrip:dev .

dev:
	docker-compose --env-file ./.env.local up

dump:
	scp donnie@airstrip:/home/donnie/airstrip/db/sql.db ./db/sql.db

sql:
	sqlite3 db/sql.db

clean:
	@echo "Cleaning Docker environment..."
	docker-compose stop
	docker-compose down -v

build:
	@echo "Building for prod"
	docker build -t donnieashok/airstrip:prod .

deploy: build
	docker push donnieashok/airstrip:prod
	@echo "Deployed!"

# Prod
live:
	ssh donnie@airstrip sudo docker pull donnieashok/airstrip:prod
	- ssh donnie@airstrip sudo docker stop airstrip
	- ssh donnie@airstrip sudo docker rm airstrip
	ssh donnie@airstrip 'mkdir -p ~/airstrip/db'
	scp ./.env donnie@airstrip:~/airstrip/
	ssh donnie@airstrip 'sudo docker run -d --restart on-failure -v ~/airstrip/db:/db --env-file ~/airstrip/.env --name airstrip donnieashok/airstrip:prod'
	ssh donnie@airstrip 'rm ~/airstrip/.env'
	@echo "Is live"
