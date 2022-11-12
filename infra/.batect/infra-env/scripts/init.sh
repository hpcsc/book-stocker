#!/bin/sh

TF_ENVIRONMENT=${1:-aws}

echo "=== Generating provider.tf for ${TF_ENVIRONMENT}"
envsubst < "provider.tf.${TF_ENVIRONMENT}" | tee provider.tf

echo "=== terraform init"
terraform init
