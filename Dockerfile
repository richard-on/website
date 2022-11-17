FROM golang:1.19.3-buster as builder

WORKDIR /website

COPY go.* ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -ldflags "-X main.version=0.1.0 -X main.build=`date -u +.%Y%m%d.%H%M%S`" \
    -o run cmd/website/main.go

FROM alpine:latest

WORKDIR /website

COPY --from=builder /website/run /website/run
# COPY --from=builder /website/.env /website/.env
# COPY --from=builder /website/etc /website/etc
COPY --from=builder /website/static /website/static

EXPOSE 80

RUN mkdir -p /website/logs && \
    apk update && apk add curl && apk add --no-cache bash && \
    apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ./run