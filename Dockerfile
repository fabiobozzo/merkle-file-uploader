FROM golang:1.21 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o ./mfu ./main.go

FROM alpine
EXPOSE 8080
COPY --from=builder /app/mfu /mfu
CMD ["/mfu", "server"]
