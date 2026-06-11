---
description: "Build the project for production"
---

Build the Go project for production.

Output is a binary in the current directory (or `go build -o <name>` for a custom output path).

```bash
go build -o tuix .
```

Cross-compile for different platforms:

```bash
GOOS=linux GOARCH=amd64 go build -o tuix-linux .
GOOS=darwin GOARCH=amd64 go build -o tuix-darwin .
GOOS=windows GOARCH=amd64 go build -o tuix.exe .
```
