# Stage 1: Base
FROM golang:1.22-alpine AS base

# Install CA certificates
RUN apk add --no-cache ca-certificates

WORKDIR /src

ARG PORT_ARG=8080
ARG SAVES_DIR_ARG=/root/backups
ARG CERT_FILE_ARG=/root/server.crt
ARG KEY_FILE_ARG=/root/server.key

COPY server.crt /src/server.crt
COPY server.key /src/server.key
COPY go.mod go.sum ./
RUN go mod download


COPY cmd/news-alligator/web /src/cmd/news-alligator/web
COPY backups/ /src/backups
COPY pkg /src/pkg

RUN go build -o /bin/web /src/cmd/news-alligator/web/main.go

# Stage 2: Run
FROM scratch
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/

# Define environment variables using build arguments
ENV PORT=${PORT_ARG}
ENV SAVES_DIR=${SAVES_DIR_ARG}
ENV CERT_FILE=${CERT_FILE_ARG}
ENV KEY_FILE=${KEY_FILE_ARG}

COPY --from=base /src/backups /root/backups
COPY server.crt /root/server.crt
COPY server.key /root/server.key

# Declare a volume
VOLUME /root/backups

COPY --from=base /bin/web /bin/
EXPOSE ${PORT}

ENTRYPOINT ["/bin/web"]