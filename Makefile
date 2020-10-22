
export GO111MODULE=on

.PHONY: bin
bin: fmt vet
	go build -o bin/updatecontext github.com/borisputerka/updatecontext/cmd/plugin

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...
