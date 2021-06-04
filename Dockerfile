## go backend
FROM golang:1.15 as go

WORKDIR /src
COPY . .
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.40.1 && \
      mv ./bin/golangci-lint /bin/.
RUN make backend-lint
RUN make backend-build

# ## shrink the binary
# ENV UPX_VER=3.96
# RUN apt update && \
#       apt install -qy xz-utils && \
#       wget https://github.com/upx/upx/releases/download/v${UPX_VER}/upx-${UPX_VER}-amd64_linux.tar.xz && \
#       tar xvf upx-${UPX_VER}-amd64_linux.tar.xz  && \
#       ./upx-${UPX_VER}-amd64_linux/upx --no-progress lander

## build front
FROM node:15.14.0 as node

WORKDIR /src
COPY . .
RUN make frontend-lint
RUN make frontend-build

## pull it all together
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=go   /src/lander .
COPY --from=node /src/frontend/dist /app/frontend/dist

## show us the sized just for the observability
RUN find /app -type d -exec du -sh {} \; ; du -sh /app/lander

ENTRYPOINT [ "./lander" ]
