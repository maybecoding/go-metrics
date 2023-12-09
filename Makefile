.PHONY: all
all: ;

.PHONY: pg
pg:
	docker run --rm \
		--name=go-metrics-db \
		-v $(abspath ./db/init/):/docker-entrypoint-initdb.d \
		-v $(abspath ./db/data/):/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD="P@ssw0rd" \
		-d \
		-p 5432:5432 \
		postgres:16.1

.PHONY: stop-pg
stop-pg:
	docker stop go-metrics-db

.PHONY: clean-data
clean-data:
	sudo rm -rf ./db/data/

.PHONY: build
build:
	go build -o cmd/server/server ./cmd/server/*.go
	go build -o cmd/agent/agent ./cmd/agent/*.go
