#!/bin/bash

set -euo pipefail

COMPONENT_NAME=${1}

gotestsum --format testname --junitfile test-result-component-${COMPONENT_NAME}.xml -- -coverprofile=coverage-component-${COMPONENT_NAME}.out -tags=component,${COMPONENT_NAME} ./...
go tool cover -html=coverage-component-${COMPONENT_NAME}.out -o coverage-component-${COMPONENT_NAME}.html
