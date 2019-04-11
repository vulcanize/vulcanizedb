BIN = $(GOPATH)/bin
BASE = $(GOPATH)/src/$(PACKAGE)
PKGS = go list ./... | grep -v "^vendor/"

default: build

GOVENDOR = $(BIN)/govendor
$(BIN)/govendor:
	go get -u github.com/kardianos/govendor

GINKGO = $(BIN)/ginkgo
$(BIN)/ginkgo:
	go get -u github.com/onsi/ginkgo/ginkgo

GX = $(BIN)/gx
$(BIN)/gx:
	go get -v -u github.com/whyrusleeping/gx
	$(GX) install --global

LINT = $(BIN)/golint
$(BIN)/golint:
	go get github.com/golang/lint/golint

METALINT = $(BIN)/gometalinter.v2
$(BIN)/gometalinter.v2:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	$(METALINT) --install

.PHONY: installtools
installtools: | $(LINT) $(GINKGO) $(GX) $(GOVENDOR)
	echo "Installing tools"

.PHONY: metalint
metalint: | $(METALINT)
	$(METALINT) ./... --vendor \
	--fast \
	--exclude="exported (function)|(var)|(method)|(type).*should have comment or be unexported" \
	--format="{{.Path.Abs}}:{{.Line}}:{{if .Col}}{{.Col}}{{end}}:{{.Severity}}: {{.Message}} ({{.Linter}})"

.PHONY: lint
lint:
	$(LINT) $$($(PKGS)) | grep -v -E "exported (function)|(var)|(method)|(type).*should have comment or be unexported"

.PHONY: test
test: | $(GINKGO) $(LINT)
	go vet ./...
	go fmt ./...
	$(GINKGO) -r

.PHONY: govendor
dep: | $(GOVENDOR)
	$(GOVENDOR) fetch -v +missing

build: govendor fixlibcrypto
	go fmt ./...
	go build

.PHONY: fixlibcrypto
fixlibcrypto:
	mkdir tmp-go-ethereum
	git clone https://github.com/ethereum/go-ethereum.git tmp-go-ethereum
	cp -r "tmp-go-ethereum/crypto/secp256k1/libsecp256k1" "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"
	rm -rf "tmp-go-ethereum"
