
PROJECTNAME := $(shell basename "$(PWD)")

#project related variables
PRJBASE := $(shell pwd)
PRJBIN  := $(PRJBASE)/bin
GOFILES := ./cmd/server/...
DB_ID   := $(shell docker ps -a -q --filter="name=ui_test")

export POSTGRESQL_URL='postgres://ui_test:postgres@localhost:5432/ui_test?sslmode=disable'

.PHONY: db-run
db-run:
	@docker run -d --name ui_test -e POSTGRES_USER=ui_test -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres

.PHONY: db-stop
db-stop:
	@docker stop $(DB_ID)
	@docker rm $(DB_ID)


.PHONY: db-migrate-up
db-migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path configs/migrations up

.PHONY: db-migrate-down
db-migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path configs/migrations down

. PHONY: docker-build
docker-build:
	@docker build --file build/Dockerfile --tag ui-assignment:v1.0.0 .

. PHONY: docker-compose-up
docker-compose-up:
	@cd deployments; docker-compose up

. PHONY: docker-compose-down
docker-compose-down:
	@cd deployments; docker-compose down


.PHONY: dev-setup
dev-setup: migrate-cli-setup
	@docker pull postgres

.PHONY: migrate-cli-setup
migrate-cli-setup:
	@sudo curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
	@sudo echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
	@sudo apt-get update
	@sudo apt-get install -y migrate
	@migrate create -ext sql -dir configs/migrations create_users_table

.PHONY: generate
generate:
	@go generate $(GOFILES)

.PHONY: clean
build: clean
	@gofmt -s -w .
	@go build -o $(PRJBIN)/$(PROJECTNAME) $(GOFILES)

.PHONY: clean
clean:
	@go clean
	@rm -rf  ${PRJBIN}/

.PHONY: run
run:
	@./bin/$(PROJECTNAME)

.PHONY: swagger-init
swagger-init:
	cd ./cmd/server; swag init --parseDependency;
	rm -rf ./api/
	mv ./cmd/server/docs ./api/
	rm -rf ./cmd/server/docs
