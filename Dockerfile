FROM ubuntu:18.04

RUN apt-get update && apt-get install -y ca-certificates

COPY bin/bot/bot .

ENTRYPOINT ./bot