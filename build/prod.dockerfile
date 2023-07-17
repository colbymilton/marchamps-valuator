FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY /pkg ./pkg
COPY /internal ./internal
COPY /cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /mv ./cmd/valuator

EXPOSE 9999

CMD ["/mv"]