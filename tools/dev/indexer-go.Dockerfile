FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY proto/go.mod proto/go.sum ./proto/
COPY services/indexers/go.mod services/indexers/go.sum ./services/indexers/
RUN cd services/indexers && go mod download

COPY proto/ proto/
COPY services/indexers/ services/indexers/
RUN cd services/indexers && go build -o /indexer-go ./cmd/indexer-go

FROM golang:1.25-alpine

RUN apk add --no-cache git ca-certificates
RUN go install github.com/sourcegraph/scip-go/cmd/scip-go@latest

COPY --from=builder /indexer-go /usr/local/bin/indexer-go

EXPOSE 9090
CMD ["indexer-go"]
