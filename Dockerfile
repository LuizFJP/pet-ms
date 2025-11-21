FROM golang:1.23 AS builder

ARG MAIN_PACKAGE=./init
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_TIME=unknown

WORKDIR /app

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

COPY . /app
COPY proto /app/proto
COPY googleapis /app/googleapis

RUN echo && ls
RUN protoc \
    -I=. \
    -I=./googleapis \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./proto/pet-ms.proto

RUN mkdir -p /app/bin && \
    go build -trimpath \
      -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
      -o /app/bin/pet-ms ${MAIN_PACKAGE}

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /app/bin/pet-ms /app

EXPOSE 50051
USER nonroot:nonroot
ENTRYPOINT ["/app"]