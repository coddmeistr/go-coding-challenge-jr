FROM golang:1.22.1-alpine3.19 as gobuild

RUN apk add -U --no-cache ca-certificates

ARG PORT=6000

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/server -o ./build

FROM scratch
WORKDIR /app
COPY --from=gobuild ./etc/ssl/certs ../etc/ssl/certs
COPY --from=gobuild ./usr/share/ca-certificates ../usr/share/ca-certificates
COPY --from=gobuild ./app/configs ./configs
COPY --from=gobuild ./app/cmd/server/build .

EXPOSE $PORT
CMD ["./build"]
