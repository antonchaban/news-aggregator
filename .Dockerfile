# Stage 1: Base
FROM golang:1.22-alpine AS base

# Install CA certificates
RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY backups /src/backups
COPY cmd/news-alligator/web /src/cmd/news-alligator/web
COPY internal /src/internal
COPY pkg /src/pkg
COPY server.crt /src/server.crt
COPY server.key /src/server.key

RUN go build -o /bin/web /src/cmd/news-alligator/web/main.go

# Stage 2: Run
FROM scratch
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/

ENV PORT=8080
ENV SAVES_DIR=/root/backups
ENV CERT_FILE=/root/server.crt
ENV KEY_FILE=/root/server.key

COPY --from=base /src/backups /root/backups
COPY server.crt /root/server.crt
COPY server.key /root/server.key

# Declare a volume
VOLUME /root/backups

COPY --from=base /bin/web /bin/
EXPOSE 8080

ENTRYPOINT ["/bin/web"]