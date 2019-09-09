FROM alpine

RUN apk add --no-cache curl

ENV CONSUL_HOST_ADDR users.consul:8500
ENV USERS_DB_KEY usersdb

ENTRYPOINT "bin/run_service.sh"

EXPOSE 443