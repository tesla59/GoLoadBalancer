FROM golang:1.21.1-alpine3.18 AS build
WORKDIR /app
COPY . .
RUN go build -o worker *.go

FROM alpine:3.14
WORKDIR /app
COPY --from=build /app/worker .
EXPOSE 8080
CMD ["./worker"]
