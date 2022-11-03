FROM golang:buster AS build

WORKDIR /build

COPY . /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o target/gobblerd ./cmd/daemon/...

FROM busybox:latest

COPY --from=build /build/target/gobblerd /usr/bin/gobblerd
COPY config.local.yml /etc/gobblerd/config.yml

CMD /usr/bin/gobblerd