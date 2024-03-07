FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

EXPOSE 6000

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/server -o /server

CMD ["/server"]