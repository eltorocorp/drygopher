

all: install test

install:
	@go install
.PHONY: install

test:
	@cd internal && mockery -all
	@drygopher -d
.PHONY: test