.PHONY: build
build:
	go build -race -o ./bin/httpd-log-monitor

.PHONY: test
test:
	go test -race -coverprofile c.out ./...

.PHONY: clean
clean:
	rm -rf ./bin ./c.out

.PHONY: lint
lint:
	golangci-lint run --tests=false
