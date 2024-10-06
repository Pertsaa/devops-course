FROM golang:1.23.0 AS build-stage
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

FROM gcr.io/distroless/base-debian12 AS build-prod-stage
WORKDIR /
COPY --from=build-stage /api /api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/api"]