FROM golang:1.17 AS builder
COPY . /var/app
WORKDIR /var/app
RUN CGO_ENABLED=0 go build -o comment-sentiment .

FROM alpine:3.14
COPY --from=builder /var/app/comment-sentiment /var/app/comment-sentiment
CMD ["/var/app/comment-sentiment"]
