FROM golang:1.21 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
