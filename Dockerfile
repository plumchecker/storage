FROM golang:latest AS builder
ENV PROJECT_PATH=/app/storage
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY . ${PROJECT_PATH}
WORKDIR ${PROJECT_PATH}
RUN go build cmd/storage/main.go

FROM golang:alpine
WORKDIR /app/cmd/storage
COPY --from=builder /app/storage/main .
EXPOSE 30001
CMD ["./main"]
