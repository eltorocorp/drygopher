

all: install test

install:
	@go install
.PHONY: install

test:
	@cd internal && mockery -all
	@drygopher -d -e "/mocks,/interfaces,/cmd,/host,'iface$$','drygopher$$'"
.PHONY: tes