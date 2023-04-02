#!/bin/bash

set -euo pipefail

COMPONENT_NAME=${1}

go test -coverprofile=coverage-component-${COMPONENT_NAME}.out -tags=component,${COMPONENT_NAME} -v ./...
go tool cover -html=coverage-component-${COMPONENT_NAME}.out -o coverage-component-${COMPONENT_NAME}.html
