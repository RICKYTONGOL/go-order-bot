all: test vet fmt lint build

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l
	test -z $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l)

install-lint:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go list -f {{.Target}} golang.org/x/lint/golint

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

install-staticcheck:
	cd /tmp && GOPROXY="" go get honnef.co/go/tools/cmd/staticcheck

staticcheck:
	staticcheck ./...

build:
	go build -o bin/api-server ./cmd/api-server
