#!/bin/bash

set -euo pipefail

go test -coverprofile=coverage-integration.out -tags=integration -v ./...
go tool cover -html=coverage-integration.out -o coverage-integration.html
