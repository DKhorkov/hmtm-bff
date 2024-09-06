## Usage

### Run via docker:

To run app and it's dependencies in docker, use next command:
```bash
make -C scripts prod
```

### Run via source files:

To run application via source files, use next commands:
```shell
go run ./cmd/hmtmbff/main.go
```

## Linters

```shell
golangci-lint run -v --fix
```

## Tests

```shell
go test -v ./test...
```
