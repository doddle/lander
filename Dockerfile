## go backend
FROM golang:1.15 as go

WORKDIR /src
COPY . .
RUN make build-go && du -sh lander

## shrink the binary
ENV UPX_VER=3.96
RUN apt update && \
      apt install -qy xz-utils && \
      wget https://github.com/upx/upx/releases/download/v${UPX_VER}/upx-${UPX_VER}-amd64_linux.tar.xz && \
      tar xvf upx-${UPX_VER}-amd64_linux.tar.xz  && \
      ./upx-${UPX_VER}-amd64_linux/upx --no-progress lander

## build front
FROM node:12.14.1 as node

WORKDIR /src
COPY . .
RUN make build-node

## pull it all together
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=go   /src/lander .
COPY --from=node /src/frontend/dist /app/frontend/dist
RUN find /app -type d -exec du -sh {} \; ; du -sh /app/lander

ENTRYPOINT [ "./lander" ]
