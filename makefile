

all: install test

install:
	@echo Updating dependencies...
	@dep ensure
	@cd drygopher && go install 
.PHONY: install

test:
	@echo Purging old mocks...
	@rm -drf drygopher/mocks
	@echo Building mocks...
	@mockery -output drygopher/mocks -dir drygopher/coverage -all 
	@echo Ready to test.
	@drygopher -d -e "/mocks,/interfaces,/cmd,/host,'iface$$','drygopher$$'"
.PHONY: test