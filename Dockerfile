FROM golang:1.12.4-stretch as build
ARG BUILD_FLAGS
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

RUN mkdir -p \
    artifacts/linux \
    artifacts/arm \
    artifacts/windows \
    artifacts/darwin
RUN GOOS=linux GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o artifacts/linux/idl ./cmd/idl
RUN GOOS=linux GOARCH=arm go build -ldflags "$BUILD_FLAGS" -o artifacts/arm/idl ./cmd/idl
RUN GOOS=darwin GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o artifacts/darwin/idl ./cmd/idl
RUN GOOS=windows GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o artifacts/windows/idl.exe ./cmd/idl

RUN GOOS=linux GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o artifacts/linux/idl-repository ./cmd/idl-repository

RUN go test ./... -race -coverprofile=coverage.txt -covermode=atomic

FROM build AS package
RUN apt-get -y update && apt-get -y install zip
RUN mkdir -p /output/artifacts
WORKDIR /build/artifacts
RUN cd linux && tar zcf /output/artifacts/linux.tar.gz idl
RUN cd arm && tar zcf /output/artifacts/arm.tar.gz idl
RUN cd darwin && tar zcf /output/artifacts/darwin.tar.gz idl
RUN cd windows && zip -r /output/artifacts/windows.zip idl.exe
WORKDIR /output
COPY --from=build /build/coverage.txt ./

FROM ubuntu:18.04 as idl
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get -y install --no-install-recommends ca-certificates 2>&1
ENTRYPOINT ["/app/idl"]
VOLUME /data
WORKDIR /data
COPY --from=0 /build/artifacts/linux/idl /app/

FROM ubuntu:18.04 as idl-repository
ENTRYPOINT ["/app/idl-repository", "-s", "/var/idl-repository"]
VOLUME /var/idl-repository
EXPOSE 80
WORKDIR /app
COPY --from=0 /build/artifacts/linux/idl-repository /app/
