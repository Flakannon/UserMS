project_name = user-service

build:
	docker build -t $(project_name):latest .

run:
	docker compose up -d

stop:
	docker compose down

fresh: stop build run