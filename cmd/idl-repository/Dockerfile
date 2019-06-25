FROM golang:1.12.4-stretch as build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o idl-repository ./cmd/idl-repository

FROM ubuntu:18.04 as final

EXPOSE 80

VOLUME /var/idl-repository

WORKDIR /app

COPY --from=0 /build/idl-repository /app/

ENTRYPOINT ["/app/idl-repository", "-s", "/var/idl-repository"]
