export GO111MODULE=on
test:
	go test -mod=readonly -v ./...
