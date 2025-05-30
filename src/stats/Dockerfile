FROM golang:1.23

RUN apt-get update && apt-get install -y \
    gcc \
    librdkafka-dev \
    pkg-config \
    git

WORKDIR /practice-security/stats
COPY go.* .

RUN go mod download -x

COPY * .

RUN go build

ENTRYPOINT ["./stats"]
