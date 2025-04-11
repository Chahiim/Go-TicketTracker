include .envrc

.PHONY: run/tests
run/tests: vet
	go test -v ./...

.PHONY: fmt
fmt: 
	go fmt ./...

.PHONY: vet
vet: fmt
	go vet ./...

.PHONY: run
run: vet
	go run ./cmd/web -addr=${ADDRESS} -dsn=${FEEDBACK_DB_DSN}


.PHONY: db/psql
db/psql:
	psql ${FEEDBACK_DB_DSN}


## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${FEEDBACK_DB_DSN} up

# db/migrations/down-1: undo the last migration
.PHONY: db/migrations/down-1
db/migrations/down-1:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${FEEDBACK_DB_DSN} down 1
