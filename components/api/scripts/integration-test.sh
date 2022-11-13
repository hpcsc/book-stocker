#!/bin/bash

set -euo pipefail

go test -coverprofile=coverage-unit.txt -tags=integration -v ./...
