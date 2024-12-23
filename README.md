## Usage

### Run via docker:

To run app and it's dependencies in docker, use next command:
```bash
task -d scripts docker_prod -v
```

### Run via source files:

To run application via source files, use next commands:
```shell
go run ./cmd/hmtmbff/main.go
```

## GraphQL

### Base files generation:

To generate all necessary for GraphQL files, use next command:

```shell
task -d scripts graphql_generate -v
```

## Linters

To run linters, use next command:

```shell
 task -d scripts linters -v
```

## Tests

To run test, use next commands. Coverage info will be
recorded to ```tests/coverage``` folder:
```shell
task -d scripts tests -v
```

## Benchmarks

To run benchmarks, use next command:

```shell
task -d scripts bench -v
```
