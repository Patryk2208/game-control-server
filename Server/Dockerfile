FROM golang:1.24-alpine AS build_stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/Server

FROM alpine:3.16 AS run_stage
RUN apk add --no-cache ca-certificates
COPY --from=build_stage /app/Server /Server

EXPOSE 80
CMD ["/Server"]