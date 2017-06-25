FROM alpine:3.5

MAINTAINER bullet@furqansoftware.com

RUN apk add --no-cache bash build-base openssh-client

ADD bullet /usr/bin/bullet
