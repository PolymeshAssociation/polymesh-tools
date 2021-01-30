ARG golangbase=1.15
FROM golang:${golangbase} as gobuild

ADD src/updater.go /opt/updater/updater.go
WORKDIR /opt/updater
RUN go get github.com/AdamSLevy/jsonrpc2 && \
    go get github.com/fsnotify/fsnotify && \
    go build

RUN chmod 0755 /opt/updater/updater

FROM gcr.io/distroless/cc

COPY --chown=4006:4006 --from=gobuild /opt/updater/updater /usr/local/bin/updater

USER 4006:4006

ENTRYPOINT ["/usr/local/bin/updater"]
CMD [ "/reserved_peers.txt" ]

