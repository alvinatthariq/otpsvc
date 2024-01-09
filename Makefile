.PHONY: run
run: run-test
	@go build . 
	@go run . 

.PHONY: run-test
run-test:
	@go test -v -tags dynamic `go list ./... | grep -i 'domain'` -cover