FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /app/
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/main ./cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin ./bin
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080

CMD ["./bin/main"]
