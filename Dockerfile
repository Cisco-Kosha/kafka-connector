FROM golang:alpine as builder
RUN apk update && \
    apk add --no-cache --virtual build-dependencies build-base gcc librdkafka librdkafka-dev pkgconf autoconf automake libtool lz4-libs zstd-dev libsasl
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o kafka-connector .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /usr/lib/pkgconfig /usr/lib/pkgconfig
COPY --from=builder /usr/lib/* /usr/lib/
COPY --from=builder /build/kafka-connector /app/
WORKDIR /app
CMD ["./kafka-connector"]