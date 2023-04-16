FROM golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0

ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY . .

COPY config/prod.yaml /app/config.yaml
RUN go build -ldflags="-s -w" -o /app/main cmd/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York

ENV TZ Europe/Moscow

WORKDIR /app

COPY --from=builder /app/* /app/
COPY --from=builder /app/config.yaml /app/config.yaml

EXPOSE 8080/tcp

CMD ["/app/main",  "--config=/app/config.yaml"]
