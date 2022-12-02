FROM golang:1.19-alpine AS build

WORKDIR /go/src/portfolioTask
COPY . .
RUN go build -o bin/server api/cmd/main.go
CMD ["./bin/server"]