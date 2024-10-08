# include variables from the .envrc file
include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@printf 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${GAZEBO_DB_DSN} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASSWORD}

## wrun/api: run the cmd/api application with auto reload
.PHONY: wrun/api
wrun/api:
	wgo make run/api

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GAZEBO_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@printf 'Creating migration files for ${name}...\n'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@printf 'Running up migrations...\n'
	migrate -path ./migrations -database ${GAZEBO_DB_DSN} up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor


# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = 'api.gazebosoftware.com'

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh gazebo@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api ubuntu@${production_host_ip}:~
	rsync -rP --delete ./migrations ubuntu@${production_host_ip}:~
	rsync -P ./remote/production/api.service ubuntu@${production_host_ip}:~
	rsync -P ./remote/production/Caddyfile ubuntu@${production_host_ip}:~
	ssh -t ubuntu@${production_host_ip} '\
		migrate -path ~/migrations -database $$GAZEBO_DB_DSN up \
		&& sudo mv ~/api.service /etc/systemd/system/ \
		&& sudo systemctl enable api \
		&& sudo systemctl restart api \
		&& sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
		'
