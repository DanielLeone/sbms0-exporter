FROM golang:1.21 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /sbms-exporter

FROM scratch
COPY --from=builder /sbms-exporter /sbms-exporter
EXPOSE 9000
ENTRYPOINT ["/sbms-exporter"]
