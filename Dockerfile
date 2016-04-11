FROM golang:1.6-alpine

RUN apk -U add git && \
  go get github.com/constabulary/gb/... && \
  git clone --depth=1 https://github.com/ubergesundheit/simple-time-series-server.git /stss-git && \
 cd /stss-git && gb build all && mkdir -p /stss && cp bin/simple-time-series-server /stss && rm -rf /stss-git /go && apk del git

WORKDIR /stss

EXPOSE 8080

ENTRYPOINT ["./simple-time-series-server"]

