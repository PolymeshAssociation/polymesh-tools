ARG golangbase=1.15
FROM golang:${golangbase} as gobuild

ADD src/ /opt/updater/
WORKDIR /opt/updater
RUN go build && \
    chmod 0755 /opt/updater/polymesh-peer-updater

FROM gcr.io/distroless/cc

COPY --chown=4006:4006 --from=gobuild /opt/updater/polymesh-peer-updater /usr/local/bin/updater

USER 4006:4006

ENTRYPOINT ["/usr/local/bin/updater"]
CMD [ "/reserved_peers.txt" ]

