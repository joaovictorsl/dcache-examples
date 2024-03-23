.PHONY: build run container

name?=dcache

build:
	go build -o bin/dcache-server

run: build
	./bin/dcache-server

container: build
	docker stop -t=0 $(name); \
	docker rm $(name); \
	docker build -t dcache-server .; \
	docker run -d --name $(name) -p $(port):3000 dcache-server; \
	echo "======CONTAINER $(name) IS RUNNING======"; \
	echo "==============$(name) LOGS=============="; \
	docker logs -f $(name)
