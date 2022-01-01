all:
	mkdir -p ./var/stock.go/log
	mkdir -p ./var/stock.go/leveldb
	cmd/version.sh
	go build -o stock.go cmd/main.go cmd/version.go

test:
	mkdir -p ./var/ut/log
	mkdir -o ./var/ut/leveldb


clean:
	@rm gostock

.PHONY: clean test