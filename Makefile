.PHONY: build

build:
	go build -ldflags "-s -w" .

clean:
	rm hd-awake