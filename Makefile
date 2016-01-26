PORT := 9020

build:
	go build -o panda

stop:
	-lsof -t -i:${PORT} | xargs kill

run: build
	nohup ./panda>/dev/null 2>&1 &

test:
	go test helper/***

.PHONY: test
