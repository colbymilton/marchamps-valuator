FROM golang:latest

WORKDIR /app
COPY . .
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -build="go build -o /build/app ./cmd/valuator" -command="/build/app" -verbose -polling="true" -polling-interval="1000" \
-exclude-dir=".git" -exclude-dir=".vscode" -exclude-dir="web" -exclude-dir="tmp" -exclude-dir="docs" -exclude-dir="build"