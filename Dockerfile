FROM golang:1.18.3-alpine3.16

RUN apk add git --no-cache

WORKDIR /app

CMD [ "git status" ]
