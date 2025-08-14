LOCAL_BIN:=$(CURDIR)/bin
APP_NAME:=moocdsl

.PHONY: run build build-run test install-golangci-lint lint validate-lint-config

# run app
run:
	./bin/$(APP_NAME)

# build
build:
	go build -o ./bin/$(APP_NAME) ./cmd

# build and run
build-run:
	go build -o ./bin/$(APP_NAME) ./cmd && ./bin/$(APP_NAME)

# run tests
test:
	go test ./... -v

# install linter
install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

# run linting
lint:
	$(LOCAL_BIN)/golangci-lint run --config=.golangci.yaml ./...

# validate golangci-lint config
validate-lint-config:
	$(LOCAL_BIN)/golangci-lint config verify --config=.golangci.yaml
