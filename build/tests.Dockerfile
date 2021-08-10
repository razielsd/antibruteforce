FROM golang:1.16-alpine

ENV GO111MODULE=on

WORKDIR /app

COPY ./ ./
RUN apk add --no-cache git make bash build-base
RUN make build

#auto stop :)
CMD ["sleep", "3600s"]