.PHONY: dependency unit-test integration-test build run cover

dependency:
	@glide install

unit-test:
# info about flags could be found by running commands "go help test" and "go help testflag"
	@go test shortener/... -cover -covermode atomic -cpu 1,2,4 -parallel 1 -short

integration-test:
	@go test shortener/... -cover -covermode atomic

integration-test-debug:
	@go test shortener/... -v -cover -covermode atomic		

build:
	@go build

install:
	@go install	

run: build
	@./shortener

cover: integration-test
	@go tool cover -html=c.out
