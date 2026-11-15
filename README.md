# nsh

My take on <https://app.codecrafters.io/courses/shell>.

## Run

```bash
go run ./app
```

## Run tests

```bash
go test -v ./tests/
go test -v ./tests/ -run 'TestPwd'
```

## Linters

Install:

```bash
curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.15.0
```

Run:

```bash
./bin/golangci-lint run
```
