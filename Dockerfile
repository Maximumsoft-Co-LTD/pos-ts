FROM golang:1.19.2-alpine3.16 AS builder
WORKDIR /src  
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app _cmd/main.go

FROM alpine:3.13  
WORKDIR /root/
COPY --from=builder /src/app .
CMD ["./app"] 