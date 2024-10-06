FROM golang:1.23.0 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api main.go

FROM ubuntu:22.04
WORKDIR /app
COPY --from=build /app/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
