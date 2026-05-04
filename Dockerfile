FROM alpine AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates go

COPY . .
RUN go mod download
RUN go build -o bin/phoenix-server

FROM scratch
WORKDIR /app

COPY --from=builder /app/bin/phoenix-server /app/server
COPY --from=builder /app/assets /app/assets
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG PORT=5000
EXPOSE $PORT

ENTRYPOINT ["/app/server"]
