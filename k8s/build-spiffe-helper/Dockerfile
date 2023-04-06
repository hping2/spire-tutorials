# Build spiffe-helper
FROM golang:1.17-alpine as helper-builder
WORKDIR /service
RUN apk add build-base git
RUN git clone -b upgrade-helper https://github.com/marcosdy/spiffe-helper.git .
RUN go build -o /service/spiffe-helper ./cmd/spiffe-helper

FROM alpine AS spiffe-helper
RUN apk --no-cache add dumb-init curl vim git
RUN apk add --no-cache --upgrade bash
COPY --from=helper-builder /service/spiffe-helper /opt/helper/spiffe-helper

WORKDIR /opt/helper/
