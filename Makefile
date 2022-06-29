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

.PHONY: migrate_testnet
migrate_testnet:
	yarn compile && bash ./migrate.bash bas-devnet-1,chapel
	node ./scripts/build-abi.js && cat ./build/addresses.json

.PHONY: all
all: install compile test
