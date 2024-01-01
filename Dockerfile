FROM  golang:1.21.5-alpine3.19 AS builder
WORKDIR /src  
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app _cmd/main.go

FROM  alpine:3.13  
WORKDIR /root/
COPY --from=builder /src/app .
CMD ["./app"] 

#docker buildx build --platform linux/arm64 -t thezeusthech/pos-ts-ms:v1.0.0 --load .
