# Stage 1: Base
FROM golang:1.22-alpine AS base

# Install CA certificates
RUN apk add --no-cache ca-certificates

WORKDIR /src

ARG PORT_ARG=443
ARG SAVES_DIR_ARG=/root/backups
ARG CERT_FILE_ARG=/root/server.crt
ARG KEY_FILE_ARG=/root/server.key

COPY pkg/server/server.crt /src/server.crt
COPY pkg/server/server.key /src/server.key
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/news-alligator/web /src/cmd/news-alligator/web
COPY backups/ /src/backups
COPY pkg/backuper /src/pkg/backuper
COPY pkg/filter /src/pkg/filter
COPY pkg/handler /src/pkg/handler
COPY pkg/model /src/pkg/model
COPY pkg/parser /src/pkg/parser
COPY pkg/scheduler /src/pkg/scheduler
COPY pkg/server /src/pkg/server
COPY pkg/service /src/pkg/service
COPY pkg/storage /src/pkg/storage

ENV PORT=${PORT_ARG}
ENV SAVES_DIR=${SAVES_DIR_ARG}
ENV CERT_FILE=${CERT_FILE_ARG}
ENV KEY_FILE=${KEY_FILE_ARG}

RUN go build -o /bin/web /src/cmd/news-alligator/web/main.go

# Stage 2: Run
FROM scratch
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/

# Define environment variables using build arguments
ENV PORT=${PORT:-443}
ENV SAVES_DIR=${SAVES_DIR:-/root/backups}
ENV CERT_FILE=${CERT_FILE:-/root/server.crt}
ENV KEY_FILE=${KEY_FILE:-/root/server.key}

COPY --from=base /src/backups /root/backups
COPY pkg/server/server.crt /root/server.crt
COPY pkg/server/server.key /root/server.key

# Declare a volume
VOLUME /root/backups

COPY --from=base /bin/web /bin/
EXPOSE ${PORT}

ENTRYPOINT ["/bin/web"]