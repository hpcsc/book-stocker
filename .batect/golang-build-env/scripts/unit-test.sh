#!/bin/bash

set -euo pipefail

gotestsum --format testname --junitfile test-result-unit.xml -- -coverprofile=coverage-unit.out -tags=unit ./...
go tool cover -html=coverage-unit.out -o coverage-unit.html
