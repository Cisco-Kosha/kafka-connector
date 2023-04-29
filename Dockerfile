FROM golang:alpine as builder
RUN apk update && \
    apk add --no-cache bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -tags musl -o kafka-connector .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /usr/lib/pkgconfig /usr/lib/pkgconfig
COPY --from=builder /usr/lib/* /usr/lib/
COPY --from=builder /build/kafka-connector /app/
WORKDIR /app
CMD ["./kafka-connector"]

# FROM golang:alpine AS build
# RUN sed -i -e 's/v[[:digit:]]\..*\//edge\//g' /etc/apk/repositories
# RUN apk upgrade --update-cache --available
# RUN apk add --no-cache \
#         gcc \
#         libc-dev \
#         librdkafka-dev=2.0.2-r0 \
#         pkgconf
# RUN mkdir /app
# WORKDIR /app
#
# ADD . /app/
# RUN go build -o kafka-connector .
#
#
# FROM alpine
# RUN sed -i -e 's/v[[:digit:]]\..*\//edge\//g' /etc/apk/repositories
# RUN apk upgrade --update-cache --available
#
# RUN apk add --no-cache \
#         librdkafka-dev=2.0.2-r0
# WORKDIR /app
# COPY --from=build /app/kafka-connector /app/
# CMD ["/app/kafka-connector"]