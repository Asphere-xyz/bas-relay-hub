.PHONY: install
install:
	yarn

.PHONY: test
test:
	yarn test

.PHONY: compile
compile:
	yarn compile && node ./scripts/build-abi.js

.PHONY: abi
abi: compile
	mkdir -p ./relayer/abigen
	abigen --abi=./build/abi/RelayHub.json --type=RelayHub --pkg=abigen --lang=go --out=./relayer/abigen/basrelayhub.go

.PHONY: all
all: install test compile
