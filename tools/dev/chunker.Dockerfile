FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src
COPY proto/go.mod proto/go.sum ./proto/
COPY services/chunkers/go.mod services/chunkers/go.sum ./services/chunkers/
RUN cd services/chunkers && go mod download

COPY proto/ proto/
COPY services/chunkers/ services/chunkers/
RUN cd services/chunkers && CGO_ENABLED=1 go build -o /chunker ./cmd/chunker

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /chunker /usr/local/bin/chunker

EXPOSE 9091
CMD ["chunker"]
