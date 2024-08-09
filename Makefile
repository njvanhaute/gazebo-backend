run:
	go run ./cmd/api

wrun:
	wgo run ./cmd/api

psql:
	psql ${GAZEBO_DB_DSN}

migration:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GAZEBO_DB_DSN} up