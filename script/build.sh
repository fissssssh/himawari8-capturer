#!/bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/himawari8_capturer_linux_amd64 cmd/cli/main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/himawari8_capturer_macos_amd64 cmd/cli/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/himawari8_capturer_windows_amd64.exe cmd/cli/main.go
