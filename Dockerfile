FROM golang:1.21 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./


# Build
RUN GOOS=linux go build -o /sbms-exporter

EXPOSE 9000

# Run
CMD ["/sbms-exporter"]
