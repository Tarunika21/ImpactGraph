
FROM golang:1.27.2-alpine AS builder

WORKDIR /src

COPY services/api/go.mod services/api/go.sum ./services/api/
COPY packages/openapi-diff/go.mod ./packages/openapi-diff/

WORKDIR /src/services/api

RUN go mod download

COPY services/api/ /src/services/api/
COPY packages/openapi-diff/ /src/packages/openapi-diff/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/impactgraph-api .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/impactgraph-api ./impactgraph-api

EXPOSE 8081

CMD ["./impactgraph-api"]
