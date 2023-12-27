DB_URL := "postgres://operator:operator@172.17.0.2:5432/metrics?sslmode=disable"

.PHONY: all
all: ;

.PHONY: pg
pg:
#	if docker network inspect my_network &> /dev/null; then \
#    	docker network rm my_network; \
#	fi
#	docker network create --subnet=192.168.1.0/24 my_network
	docker run --rm \
		--name=go-metrics-db \
		-v $(abspath ./db/init/):/docker-entrypoint-initdb.d \
		-v $(abspath ./db/data/):/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD="postgres" \
		-d \
		-p 5432:5432 \
		postgres:16.1
#		--network=my_network --ip=192.168.1.2 \

.PHONY: pg-stop
pg-stop:
	docker stop go-metrics-db
.PHONY: pg-reset
pg-reset:
	rm -rf ./db/data/

.PHONY: build
build:
	go build -o cmd/server/server ./cmd/server/*.go
	go build -o cmd/agent/agent ./cmd/agent/*.go


.PHONY: mg-create
mg-create:
	docker run --rm \
	  -v $(realpath ./db/migrations):/migrations \
	  migrate/migrate:v4.16.2 \
	  create \
	  -dir /migrations \
	  -ext .sql \
	  -seq -digits 3 \
	  create_tables
	sudo chown -R $(whoami):staff ./db/migrations


.PHONY: mg-up
mg-up:
	docker run --rm \
	  -v $(realpath ./db/migrations):/migrations \
	  migrate/migrate:v4.16.2 \
	  -path=/migrations \
	  -database $(DB_URL) \
	  up

.PHONY: mg-down
mg-down:
	docker run --rm \
	  -v $(realpath ./db/migrations):/migrations \
	  migrate/migrate:v4.16.2 \
	  -path=/migrations \
	  -database $(DB_URL) \
	  down -all

.PHONY: test-100k
test-100k:
	curl -X POST \
	  -H "Content-Type: application/json" \
	  -H "Content-Encoding: gzip" \
	  --data-binary "@100k.json.gz" \
	  --compressed \
	  http://localhost:8080/updates/

.PHONY: test-100kh
test-100kh:
	curl -X POST \
 	  -H "Content-Type: application/json" \
 	  -H "Content-Encoding: gzip" \
 	  -H "HashSHA256: 1e53a7c6f744b1a520eeee986912db8c7af7dfd68d505ea15bb5e23cf4dbb550" \
 	  --data-binary "@100k.json.gz" \
 	  --compressed \
 	  http://localhost:8080/updates/


