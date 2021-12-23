FROM golang:1.17.3 AS build-env

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

ADD . /dockerdev
WORKDIR /dockerdev

RUN go build -gcflags '-N -l' -o /deploy_server

## Final stage
FROM debian:buster as prod

EXPOSE 9105

WORKDIR /app

COPY --from=build-env /deploy_server /app/deploy_server
COPY --from=build-env /dockerdev/config /app/config
RUN chmod +x /app/deploy_server

CMD ["/app/deploy_server"]