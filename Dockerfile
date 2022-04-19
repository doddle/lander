## go backend
FROM docker.io/library/golang:1.17 as go


WORKDIR /temp
## install goimports
RUN go mod init temp
RUN go get golang.org/x/tools/cmd/goimports@latest && go install golang.org/x/tools/cmd/goimports
# install golangci-lint
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.40.1 && \
      mv ./bin/golangci-lint /bin/.

WORKDIR /src
COPY . .
RUN find . -type f
RUN make backend-lint
RUN make backend-build


#### shrink the lander golang binary

######## shrink the binaries
FROM docker.io/hairyhenderson/upx:3.94 as upx
WORKDIR /src
COPY --from=go /src/lander .
RUN upx lander



## build front
FROM docker.io/library/node:16.13.2 as node

# This step produces /src/frontend/dist which we will COPY into the final docker image

WORKDIR /src
COPY . .
RUN make frontend-lint
RUN make frontend-build

## pull it all together
FROM docker.io/library/alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app

# Copy in the lander golang binary
COPY --from=upx   /src/lander .
# Copy in the "dist" folder from the node stage
COPY --from=node /src/frontend/dist /app/frontend/dist

## show us the size just for the observability
RUN find /app -type d -exec du -sh {} \; ; du -sh /app/lander


## create a non-root user
ENV USER=app
ENV UID=1000
ENV GID=1000

# add the group+user called $USER
# NOTE $HOME is being set to /tmp so your shell history isn't stored in the docroot (in case u kubectl exec into the pod)
RUN addgroup -g "$GID" "$USER" && \
    adduser \
      --disabled-password \
      --gecos "" \
      --home "/tmp" \
      --ingroup "$USER" \
      --uid "$UID" \
      "$USER"

# chown the files to be owned by $USER
RUN chown -R $UID:$GID /app


# set the container to run as $USER henceforth
USER $USER


# launch lander on start
ENTRYPOINT [ "./lander" ]
