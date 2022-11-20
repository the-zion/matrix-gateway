FROM debian:stable-slim

COPY ./bin/gateway /home/app/

WORKDIR /home/app/

EXPOSE 8080
EXPOSE 7070
