FROM golang:1.20-alpine

WORKDIR /practice-security/gateway
COPY go.* .

RUN go mod download -x

COPY * .

RUN go build .

ENTRYPOINT ["./gateway"]
