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
