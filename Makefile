help:			## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

all: test 		## Run only the test for now.

test: 			## Run all unit tests.
	@cd godoghttp ; go test -v -cover

test-coverage:		## Run all unit tests with coverage report.
	@cd godoghttp ; go test -coverprofile=cover.out && go tool cover -html=cover.out

lint:			## Run lint test on the package..includes gofmt, go vet and golint
	@cd godoghttp ; gofmt -l -e -s . ; go vet . ; golint .
