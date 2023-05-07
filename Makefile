.PHONY: test
test:
	go clean -testcache
	go test -v -cover -race ./...


.PHONY: build
build:
	go build  -o ./build/main github.com/Mutusva/monzo-webcrawler/cmd

