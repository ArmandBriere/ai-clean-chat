FROM golang:1.23 AS builder

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go install github.com/air-verse/air@latest

FROM golang:1.23

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    libopus-dev \
    libopusfile-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /app .

CMD ["air"]
