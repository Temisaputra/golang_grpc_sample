# Step 1: Modules caching
FROM golang:1.23 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.23 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go mod vendor

# build binary utama (namanya be-example biar sama kayak cobra command)
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -trimpath -a -tags "osusergo,netgo,migrate" \
    -ldflags '-extldflags=-static -w -s -v' \
    -o /bin/be-example ./main.go

# Step 3: Runtime
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Jakarta
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# copy hasil build
COPY --from=builder /bin/be-example /app/

# default command → jalanin gRPC server
CMD ["/app/be-example", "grpc"]
