all:
	mkdir -p ./var/stock.go/log
	mkdir -p ./var/stock.go/leveldb
	cmd/version.sh
	go build -o stock.go cmd/main.go cmd/version.go

test:
	mkdir -p ./var/ut/log
	mkdir -o ./var/ut/leveldb

pb: model/pb/*.proto
	protoc --go_out=./model/pb/ ./model/pb/*.proto

clean:
	@rm gostock

.PHONY: clean test