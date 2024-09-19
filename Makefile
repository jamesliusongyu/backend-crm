run:
	DATABASE_URL=mongodb://root:Password123@localhost:27017 go run cmd/main.go
test:
	find . -name go.mod -execdir go test ./... \;
