ARG golangbase=1.15
FROM golang:${golangbase} as gobuild

ADD src/ /opt/health-check/
WORKDIR /opt/health-check
RUN ls -a && \
    go build && \
    chmod 0755 /opt/health-check/polymesh-health-check

FROM debian:buster as releasefetch

ADD src/fetch.sh /opt/fetch/fetch.sh
WORKDIR /opt/fetch
RUN bash /opt/fetch/fetch.sh

FROM gcr.io/distroless/cc
USER 4001:4001

COPY --chown=4001:4001 --from=releasefetch /opt/files/                    /
COPY --chown=4002:4002 --from=gobuild      /opt/health-check/polymesh-health-check /usr/local/bin/check

ENTRYPOINT ["/usr/local/bin/polymesh"]
CMD [ "-d", "/var/lib/polymesh" ]

HEALTHCHECK \
    --interval=10s \
    --start-period=120s \
    --timeout=5s \
    --retries=6 \
    CMD /usr/local/bin/check liveness

