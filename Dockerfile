ARG GO_VERSION=1.16
FROM golang:$GO_VERSION-alpine AS builder

ADD cmd/server /src/cmd/server
ADD pkg/db /src/pkg/db
ADD go.mod /src/go.mod
ADD go.sum /src/go.sum
WORKDIR /src
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /server ./cmd/server

FROM alpine:latest
COPY --from=builder /src/pkg/db /pkg/db
COPY --from=builder /server /server
CMD ["./server", "-d"]