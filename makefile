

all: install test

install:
	@go install
.PHONY: install

test:
	@drygopher -d
.PHONY: test