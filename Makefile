
all:
	cmd/version.sh
	go build -o gostock cmd/main.go cmd/version.go


run: all
	mkdir -p ./var/stock.go/log
	mkdir -p ./var/stock.go/leveldb
	./gostock -c ./cmd/gostock.yaml

pb: db/leveldb/leveldb_model.proto api/api_model.proto
	protoc --go_out=./api/ ./api/api_model.proto
	protoc --go_out=./db/leveldb/ ./db/leveldb/leveldb_model.proto

clean:
	@rm gostock

.PHONY: clean
