FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /app ./cmd/urlshortener/main.go

FROM alpine:3.18
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
