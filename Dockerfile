FROM golang:1.19.3-buster as builder

WORKDIR /website

COPY go.* ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o run cmd/website/main.go

FROM alpine:3.15.4

WORKDIR /website

COPY --from=builder /website/run /website/run
COPY --from=builder /website/.env /website/.env
COPY --from=builder /website/etc /website/etc
COPY --from=builder /website/static /website/static

RUN mkdir -p /website/logs

EXPOSE 80
EXPOSE 443

RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ./run win-dev