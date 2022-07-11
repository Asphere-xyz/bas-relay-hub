install: ./package.json ./yarn.lock
	yarn

test: ./test
	yarn test

compile: ./contracts
	yarn compile && node ./scripts/build-abi.js && cat ./build/addresses.json

.PHONY: abi
abi: compile
	mkdir -p ./relayer/abigen
	abigen --abi=./build/abi/RelayHub.json --type=RelayHub --pkg=abigen --lang=go --out=./relayer/abigen/basrelayhub.go

.PHONE: proto
proto:
	protoc -I./ -I/usr/local/include --go_out=./relayer/proto --go-grpc_out=./relayer/proto ./relayer/proto/*.proto

.PHONY: migrate_testnet
migrate_testnet:
	yarn compile && bash ./migrate.bash bas-devnet-1,chapel
	node ./scripts/build-abi.js && cat ./build/addresses.json

.PHONY: all
all: install compile test
