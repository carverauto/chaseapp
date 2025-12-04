PACKAGE=chaseapp-get-streams
.DEFAULT_GOAL := darwin-arm64

GOVERSION := 1.19.3

linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(PACKAGE)-linux-arm64 main.go structs.go

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(PACKAGE)-linux-amd64 main.go structs.go

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(PACKAGE)-darwin-arm64 main.go structs.go

clean:
	rm -f $(PACKAGE)-linux-arm64
	rm -f $(PACKAGE)-linux-amd64
	rm -f $(PACKAGE)-darwin-arm64