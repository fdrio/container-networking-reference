# Builder

FROM golang:alpine as builder
WORKDIR /project
COPY main.go .
COPY go.mod .

RUN go build -o printip main.go


FROM alpine:latest

COPY --from=builder /project/printip /

ENTRYPOINT ["/printip"]

