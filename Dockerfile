FROM golang:1.21 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o operator .

FROM alpine:3.14
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/operator .
COPY --from=builder /app/templates ./templates

# Health check порт
EXPOSE 8080

# Стандартные проверки для k8s
EXPOSE 8081

CMD ["./operator"]
