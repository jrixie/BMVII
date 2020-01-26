FROM golang:1.13 as builder

ENV GO111MODULE=on

COPY . /build

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /build .

CMD ["./server"]
