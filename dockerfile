FROM golang:1.22.3-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o loadtest ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /app/loadtest .

ENTRYPOINT ["/app/loadtest"]