FROM golang:1.19.1 AS builder
WORKDIR /build

# Precompile the entire go standard library
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go install -v -a std

# Download and precompile all third party libraries, ignoring errors
ADD go.mod .
ADD go.sum .
RUN go mod download -x
RUN go list -f '{{.Path}}/...' -m all | GOARCH=amd64 GOOS=linux CGO_ENABLED=0 xargs -n1 go build -v -i; echo done

#  Add the sources
ADD . .

ARG VERSION=main

# Compile only our sources
RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
      go build -v \
              --ldflags "-extldflags -static -X 'github.com/hpcsc/book-stocker/api/internal/handler.Version=${VERSION}'" \
              -o api \
        ./cmd/api/main.go

FROM scratch
WORKDIR /app
ENTRYPOINT ["/app/api"]
COPY --from=builder /build/api .
