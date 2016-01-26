PORT := 9020

build:
	go build -o panda

stop:
	-lsof -t -i:${PORT} | xargs kill

run: build
	./panda

test:
	go test helper/***

.PHONY: test
