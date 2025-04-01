FROM golang:1.23-alpine

WORKDIR /practice-security/grpc-gateway
COPY go.* .

RUN go mod download -x

COPY * .

RUN go build

ENTRYPOINT ["./grpc-gateway"]
