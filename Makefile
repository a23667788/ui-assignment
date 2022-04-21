
PROJECTNAME := $(shell basename "$(PWD)")

#project related variables
PRJBASE := $(shell pwd)
PRJBIN  := $(PRJBASE)/bin
GOFILES := ./cmd/server/...

export POSTGRESQL_URL='postgres://ui_test:postgres@localhost:5432/ui_test?sslmode=disable'

.PHONY: run-db
db-run:
	@docker run -d --name ui_test -e POSTGRES_USER=ui_test -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres

.PHONY: db-migrate-up
db-migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path assets/configs/migrations up

.PHONY: db-migrate-down
db-migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path assets/configs/migrations down


.PHONY: dev-setup
dev-setup: migrate-cli-setup
	@docker pull postgres
	@go get -u github.com/gorilla/mux
	@go get -u gorm.io/gorm


.PHONY: migrate-cli-setup
migrate-cli-setup:
	sudo curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
	sudo echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
	sudo apt-get update
	sudo apt-get install -y migrate

.PHONY: generate
generate:
	@go generate $(GOFILES)

.PHONY: clean
build: clean
	@gofmt -s -w .
	@go build -o $(PRJBIN)/$(PROJECTNAME) $(GOFILES)

.PHONY: clean
clean:
	go clean
	rm -rf  ${PRJBIN}/