FROM golang:1.23 AS builder

ARG MAIN_PACKAGE=./init
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_TIME=unknown

WORKDIR /

RUN apt-get update && apt-get install -y --no-install-recommends \
      protobuf-compiler ca-certificates git unzip \
    && rm -rf /var/lib/apt/lists/*

ENV GOBIN=/usr/local/bin
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1 \
 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN protoc -I ./proto \
    --go_out=./proto --go_opt=paths=source_relative \
    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
    ./proto/pet-ms.proto

RUN go build -trimpath \
    -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
    -o /app ${MAIN_PACKAGE}

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /app /app

EXPOSE 50051 2112
USER nonroot:nonroot
ENTRYPOINT ["/app"]