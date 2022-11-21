FROM debian:stable-slim

COPY ./bin/gateway /home/app/
RUN useradd matrix && chown -R matrix:matrix /home/app/ && chmod 700 /home/app/
WORKDIR /home/app/
USER matrix
EXPOSE 8080
EXPOSE 7070
