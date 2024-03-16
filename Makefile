
.PHONY: test
test:
	go test ./... -race --cover -coverprofile=coverage.out

.PHONY: db.migrate
db.migrate:
	go run db/main.go
	

.PHONY: api.start
api.start:
	go run cmd/main.go
	
