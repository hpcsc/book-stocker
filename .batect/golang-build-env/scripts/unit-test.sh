#!/bin/bash

set -euo pipefail

go test -coverprofile=coverage-unit.out -tags=unit -v ./...
go tool cover -html=coverage-unit.out -o coverage-unit.html
