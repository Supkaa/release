FROM golang:1.23-alpine3.21 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release main.go

FROM alpine:3.21
RUN apk add --no-cache git
COPY --from=build /app/release /release
ENTRYPOINT ["/release"]