FROM golang:1.17-alpine as builder
RUN apk add build-base git

FROM alpine AS image-base 
RUN apk --no-cache add dumb-init curl vim git

#WORKDIR /opt/service/

# Build spiffe-helper
FROM builder as helper-builder
WORKDIR /service
RUN git clone -b upgrade-helper https://github.com/marcosdy/spiffe-helper.git .
RUN go build -o /service/spiffe-helper ./cmd/spiffe-helper

FROM image-base AS spiffe-helper
#RUN addgroup -g 70 ssl-cert
#RUN adduser -G ssl-cert -u 70 -D postgres

RUN apk add --no-cache --upgrade bash
COPY --from=helper-builder /service/spiffe-helper /opt/helper/spiffe-helper

WORKDIR /opt/helper/
