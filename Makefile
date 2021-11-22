export GO111MODULE=on

test:
	go test ./... -race -cover

run:
	go run main.go