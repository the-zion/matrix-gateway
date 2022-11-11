FROM debian:stable-slim

COPY ./bin/gateway /app/

WORKDIR /app/

EXPOSE 8080
EXPOSE 7070
