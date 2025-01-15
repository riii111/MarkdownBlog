down:
	docker compose down -v

network:
	docker network create devtrackr-network

prune:
	docker system prune --force

build:
	docker compose build

build-n:
	docker compose build --no-cache

up:
	docker compose up -d

re:
	make down
	make build
	make up

re-n:
	make down
	make build-n
	make up

re-all-force:
	make down
	make prune
	make network
	make build-n
	make up

container=frontend
command=/bin/bash
tail=200

logs:
	docker compose logs --tail=${tail} ${container}

logs-all:
	docker compose logs --tail=${tail}

exec:
	docker compose exec ${container} ${command}