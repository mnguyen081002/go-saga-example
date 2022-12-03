FROM golang:1.19-alpine as build-env

# cache dependencies first
WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY . /app

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./main.go

FROM alpine:latest

WORKDIR /app
COPY ./config/config.release.yml /app/config/config.release.yml
COPY ./config/config.dev.yml /app/config/config.dev.yml
COPY ./config/config.yml /app/config/config.yml

COPY --from=build-env /app/main /app/main

ENTRYPOINT ["/app/main"]


