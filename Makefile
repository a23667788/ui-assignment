

#project related variables
PRJBASE := $(shell pwd)
PRJBIN  := $(PRJBASE)/bin
GOFILES := ./cmd/server/...

.PHONY: run-db
db-run:
	@docker run -d --name ui_test -e POSTGRES_USER=ui_test -e POSTGRES_PASSWORD=postgres -p 8080:8080 postgres


.PHONY: dev-setup
dev-setup:
	@docker pull postgres
	@go get -u github.com/gorilla/mux
	@go get -u gorm.io/gorm



.PHONY: generate
generate:
	@go generate $(GOFILES)

.PHONY: clean
build: clean
	@go build -o $(PRJBIN)/$(PROJECTNAME) $(LDFLAGS) $(GOFILES)
	# @cp -r ${PRJBASE}/assets/* ${PRJBIN}