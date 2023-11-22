# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.20-buster AS build

WORKDIR /app

COPY . .

RUN go mod tidy
RUN make build

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build ./app/out/app /app
ENV GIN_MODE=release

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app"]