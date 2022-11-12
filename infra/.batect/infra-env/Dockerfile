FROM ghcr.io/terraform-linters/tflint:v0.40.0 AS tflint

FROM hashicorp/terraform:1.2.9

RUN apk add --no-cache gettext
COPY --from=tflint /usr/local/bin/tflint /usr/local/bin/tflint

ADD scripts /scripts
