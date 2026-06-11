---
description: "Run linters and static analysis"
---

Run static analysis on the Go codebase:

```bash
go vet ./...
```

Run golangci-lint (if configured):

```bash
golangci-lint run ./...
```

Check formatting:

```bash
gofmt -d ./
# or
test -z "$(gofmt -l ./)"
```
