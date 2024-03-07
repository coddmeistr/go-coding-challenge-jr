FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
COPY . ./
RUN go mod download

EXPOSE 6000

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/server -o ./build

CMD ["./cmd/server/build"]