#~/terraform-providers/terraform-provider-securden

HOSTNAME=terraform.local
NAMESPACE=local
NAME=securden
BINARY=terraform-provider-securden
VERSION=0.0.1
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}_v${VERSION}

install: build
	mkdir -p "${HOME}/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}"
	cp ${BINARY}_v${VERSION} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/

tools:
	#go generate -tags tools tools/tools.go
	cd tools; go generate ./...

lint:
	bin/golangci-lint run

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

tools_install:
	export GOBIN=$PWD/bin
	export PATH=$GOBIN:$PATH
	#go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	which tfplugindocs
