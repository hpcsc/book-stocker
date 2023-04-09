#!/bin/bash

set -euo pipefail

gotestsum --format testname --junitfile test-result-integration.xml -- -coverprofile=coverage-integration.out -tags=integration ./...
go tool cover -html=coverage-integration.out -o coverage-integration.html
