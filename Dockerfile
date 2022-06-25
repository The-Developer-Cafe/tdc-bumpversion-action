FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /app

COPY ./ ./

RUN go build -o bumpversion main.go

FROM alpine:3.16

RUN apk add git --no-cache

# WORKDIR /github/workspace

COPY --from=builder /app/bumpversion /bin/

ENTRYPOINT [ "bash", "-c", "pwd && ls && bumpversion" ]