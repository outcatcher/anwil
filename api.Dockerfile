FROM golang:1.20-alpine3.17 AS builder

WORKDIR /opt/build

COPY . .

RUN go build -trimpath -o ./anwil ./domains/api/cmd/server/main.go

FROM alpine:3.17.2

WORKDIR /opt/anwil

RUN apk --update --no-cache add curl # for health checks

HEALTHCHECK \
  --retries=3 \
  --interval=1m \
  --timeout=2s \
  --start-period=5s \
  CMD curl http://localhost:8010/api/v1/echo || exit 1

COPY --from=builder /opt/build/anwil ./anwil
COPY static ./static
COPY ./anwil-config.yaml ./config.yaml
COPY ./.keys ./.keys

CMD ["./anwil", "-config", "./config.yaml"]
