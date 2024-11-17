FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service-a

FROM scratch
WORKDIR /app
COPY --from=builder /app/service-a .
ENTRYPOINT ["./service-a"]