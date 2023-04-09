FROM golang:1.19.1
WORKDIR /app

ARG GOLANGCI_VERSION=1.49.0
ARG GOSEC_VERSION=2.13.1
ARG GOTESTSUM_VERSION=1.10.0
ARG TERRAFORM_VERSION=1.4.4

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_VERSION}
RUN curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOSEC_VERSION}
RUN curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz | tar xvzf - -C /usr/local/bin
RUN apt-get update && apt-get install -y unzip
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip ./terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    mv ./terraform /usr/local/bin
ADD ./scripts /build/scripts
