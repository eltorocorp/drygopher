

local: build test

prebuild:
	@echo Preparing build tooling...
	@go get -u github.com/vektra/mockery/.../
.PHONY: prebuild

build:
	@echo Updating dependencies...
	@cd drygopher && go install 
.PHONY: build

test:
	@echo Purging old mocks...
	@rm -drf drygopher/mocks
	@echo Building mocks...
	@mockery -output drygopher/mocks -dir drygopher/coverage -all 
	@echo Ready to test.
	@drygopher -d -e "/mocks,/interfaces,/cmd,/host,'iface$$','drygopher$$','types$$'"
.PHONY: test