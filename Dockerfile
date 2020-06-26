FROM golang:alpine as builder

RUN mkdir heartbeat
ADD . /heartbeat
WORKDIR /heartbeat

RUN go mod download 
RUN go build -o heartbeat ./cmd/heartbeat/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /heartbeat/heartbeat /bin/heartbeat
ENTRYPOINT heartbeat