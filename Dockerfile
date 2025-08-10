FROM alpine:3.15

LABEL maintainer="bullet@furqansoftware.com"

RUN apk add --no-cache bash build-base git make openssh-client

ADD bullet /usr/bin/bullet
