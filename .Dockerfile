FROM golang:1.22-alpine AS base
# Install CA certificates
RUN apk add --no-cache ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

#

FROM base AS build
RUN go build -o /bin/web ./cmd/news-alligator/web/main.go

#

FROM scratch
# Install CA certificates
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/
ENV PORT=8080
ENV SAVES_DIR=/root/backups
ENV CERT_FILE=/root/server.crt
ENV KEY_FILE=/root/server.key

COPY --from=build /src/backups /root/backups
COPY server.crt /root/server.crt
COPY server.key /root/server.key
COPY --from=build /bin/web /bin/

EXPOSE 8080
ENTRYPOINT ["/bin/web"]
