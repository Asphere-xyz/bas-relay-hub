.PHONY: install
install:
	yarn

.PHONY: compile
compile:
	yarn compile && node ./scripts/build-abi.js

.PHONY: all
all: install compile
