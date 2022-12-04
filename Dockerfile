FROM golang:1.19-alpine as build-env

# cache dependencies first
WORKDIR /app
COPY /dtm/go.mod /app
COPY /dtm/go.sum /app
RUN go mod download
RUN ls

COPY ./dtm /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./main.go
ENTRYPOINT ["/app/main"]