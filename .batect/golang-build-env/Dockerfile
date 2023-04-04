FROM golang:1.19.1
WORKDIR /app

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.49.0
RUN curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.13.1
RUN apt-get update && apt-get install -y unzip
RUN wget https://releases.hashicorp.com/terraform/1.4.4/terraform_1.4.4_linux_amd64.zip && \
    unzip ./terraform_1.4.4_linux_amd64.zip && \
    mv ./terraform /usr/local/bin
ADD ./scripts /build/scripts
