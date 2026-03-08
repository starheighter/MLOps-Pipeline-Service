FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./app"]