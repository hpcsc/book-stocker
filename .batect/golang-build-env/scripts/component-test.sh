#!/bin/bash

set -euo pipefail

go test -coverprofile=coverage-component.out -tags=component -v ./...
go tool cover -html=coverage-component.out -o coverage-component.html
