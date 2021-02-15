FROM golang:alpine AS builder
WORKDIR /build
COPY main.go .

RUN go build -o printIfaces

FROM alpine:latest 
COPY --from=builder /build/printIfaces /
ENTRYPOINT ["/printIfaces"]




