.PHONY: build image clean

build:
	go build -o gateway main.go

image:
	docker build -t gateway:v0.1 .

clean:
	rm -f ./gateway