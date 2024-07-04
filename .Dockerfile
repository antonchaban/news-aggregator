# Stage 1: Base
FROM golang:1.22-alpine AS base
# Install CA certificates
RUN apk add --no-cache ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/web ./cmd/news-alligator/web/main.go

# Stage 2: Run
FROM scratch
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/

ENV PORT=8080, SAVES_DIR=/root/backups, CERT_FILE=/root/server.crt, KEY_FILE=/root/server.key

COPY --from=base /src/backups /root/backups
COPY server.crt /root/server.crt
COPY server.key /root/server.key

# Declare a volume
VOLUME /root/backups

COPY --from=base /bin/web /bin/
EXPOSE 8080

ENTRYPOINT ["/bin/web"]
