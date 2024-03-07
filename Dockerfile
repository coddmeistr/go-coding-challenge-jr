FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . ./

EXPOSE 6000

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/server -o ./build

CMD ["./cmd/server/build"]