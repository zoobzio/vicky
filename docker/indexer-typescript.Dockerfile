FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY proto/go.mod proto/go.sum ./proto/
COPY indexers/go.mod indexers/go.sum ./indexers/
RUN cd indexers && go mod download

COPY proto/ proto/
COPY indexers/ indexers/
RUN cd indexers && go build -o /indexer-typescript ./cmd/indexer-typescript

FROM node:22-alpine

RUN npm install -g @sourcegraph/scip-typescript

COPY --from=builder /indexer-typescript /usr/local/bin/indexer-typescript

EXPOSE 9090
CMD ["indexer-typescript"]
