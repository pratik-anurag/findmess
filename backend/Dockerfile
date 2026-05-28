FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/findmesh-api ./cmd/api

FROM alpine:3.20
RUN adduser -D -H findmesh
USER findmesh
COPY --from=build /out/findmesh-api /usr/local/bin/findmesh-api
EXPOSE 8080
ENTRYPOINT ["findmesh-api"]
