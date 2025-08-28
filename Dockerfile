
FROM golang:latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o srv cmd/server/main.go 

EXPOSE 8080

CMD ["/srv"]
