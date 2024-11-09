export PATH := $(PATH):`go env GOPATH`/bin
export GO111MODULE=on
LDFLAGS := -s -w

all: env fmt build

build: frps

env:
	@go version

# compile assets into binary file
file:
	rm -rf ./assets/frps/static/*
	cp -rf ./web/frps/dist/* ./assets/frps/static

fmt:
	go fmt ./...

fmt-more:
	gofumpt -l -w .

gci:
	gci write -s standard -s default -s "prefix(github.com/kirilngusi/go-reverse-proxy/)" ./

vet:
	go vet ./...

frps:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags frps -o bin/frps ./cmd/frps

clean:
	rm -f ./bin/frps
