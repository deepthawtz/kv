package = github.com/deepthawtz/kv

.PHONY: install release test

install:
	go get -t -v ./...

release:
	mkdir -p release
	GOOS=darwin GOARCH=amd64 go build -o release/kv-darwin-amd64 $(package)
	GOOS=darwin GOARCH=386 go build -o release/kv-darwin-386 $(package)
	GOOS=linux GOARCH=amd64 go build -o release/kv-linux-amd64 $(package)
	GOOS=linux GOARCH=386 go build -o release/kv-linux-386 $(package)

test:
	go test -v -cover ./...
