# Build stage
FROM golang:1.25.0 AS build


WORKDIR /app
COPY . /app/
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/main.go

# Final Image stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/api ./

EXPOSE 3000

CMD ["./api"]