FROM alpine:3.7

RUN \
  apk add --update go git make gcc musl-dev linux-headers ca-certificates && \
  git clone --depth 1 https://github.com/ovcharovvladimir/essentiaHybrid && \
  (cd essentiaHybrid && make gess) && \
  cp essentiaHybrid/build/bin/gess /gess && \
  apk del go git make gcc musl-dev linux-headers && \
  rm -rf /essentiaHybrid && rm -rf /var/cache/apk/*

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/gess"]
