FROM golang:1.20

RUN apt-get update && apt-get install -y \
    gcc \
    librdkafka-dev \
    pkg-config \
    git

WORKDIR /practice-security/auth
COPY go.* .

RUN go mod download -x

COPY * .

RUN go build .
RUN go test -v ./...

ENTRYPOINT ["./auth"]
