ARG GOLANG_VERSION=1.21.0-alpine3.18

FROM golang:${GOLANG_VERSION} AS build
WORKDIR /build
COPY . .

RUN apk add git openssh wget

RUN go mod vendor

RUN go build -o /bin/api-service -mod=vendor main.go

FROM alpine:latest AS dev

WORKDIR /

EXPOSE 3000

COPY --from=build /bin/api-service /bin/api-service

ENTRYPOINT ["/bin/api-service", "--config", "config.json"]