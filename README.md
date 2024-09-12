## Usage

### Run via docker:

To run app and it's dependencies in docker, use next command:
```bash
task -d scripts docker_prod
```

### Run via source files:

To run application via source files, use next commands:
```shell
go run ./cmd/hmtmbff/main.go
```

## GraphQL

### Base files generation:
```shell
task -d scripts graphql_generate
```

## Linters

```shell
 task -d scripts linters
```

## Tests

```shell
task -d scripts tests
```
