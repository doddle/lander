########## build go binary
FROM golang:1.14 as go
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lander .


##########
FROM node:12.14.1 as node

WORKDIR /src
COPY . .
RUN cd frontend && npm i --from-lockfile && npm run build
RUN find /src/frontend/dist -type f


##########  install results in a container
FROM alpine
WORKDIR /app

COPY --from=go   /src/lander .
COPY --from=node /src/frontend/dist /app/frontend/dist

RUN find /app -type f

ENTRYPOINT ["/app/lander"]
