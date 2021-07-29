FROM golang:1.16-alpine

EXPOSE 8080
ENV GO111MODULE=on

WORKDIR /app

COPY ./ ./
RUN apk add --no-cache git make bash
RUN make build

FROM alpine

WORKDIR /app

COPY --from=0 /app/bin/abf /app/abf

CMD ["/app/abf", "server"]