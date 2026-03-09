FROM alpine AS builder

FROM scratch

ARG SERVICE=phoenix-server

COPY ./bin/${SERVICE} /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG PORT=5000
EXPOSE $PORT

ENTRYPOINT ["/server"]
