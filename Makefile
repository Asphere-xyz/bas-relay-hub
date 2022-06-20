.PHONY: install
install:
	yarn

.PHONY: test
test:
	yarn test

.PHONY: compile
compile:
	yarn compile && node ./scripts/build-abi.js

.PHONY: all
all: install test compile
