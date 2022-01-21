export GO111MODULE=on

name = wallet
port = 8080
secret = 'secret'

test:
	go test ./... -race -cover

build-image:
	docker build -t $(name) .

run:
	go run main.go

run-docker:
	docker run -it --rm --name wallet -p $(port):$(port) wallet --port=$(port) --secret=$(secret)