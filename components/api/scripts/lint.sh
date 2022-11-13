#!/bin/bash

set -euo pipefail

golangci-lint version
echo Running go vet
go vet ./...

echo Running linting
golangci-lint run
