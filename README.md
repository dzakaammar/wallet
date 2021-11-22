# wallet
Simple wallet implementation in Go

# Commands

Run the test:
```
go test ./... -race -cover
```

Run the app with default config:
```
go run main.go
```

or you can specify the http server's port and the JWT secret with:
```
go run main.go -p 8080 -s "secret"
```
