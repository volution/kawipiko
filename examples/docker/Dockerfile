FROM golang:1.18 AS kawipiko-builder

RUN mkdir -p -- /tmp/kawipiko/bin /tmp/kawipiko/src /tmp/kawipiko/go && \
  git clone -b development https://github.com/volution/kawipiko.git /tmp/kawipiko/src

WORKDIR /tmp/kawipiko/src/sources
ENV GOPATH=/tmp/kawipiko/go
RUN go build -tags 'netgo' -gcflags 'all=-l=4' -ldflags 'all=-s' -trimpath -o /tmp/kawipiko/bin/kawipiko-server ./cmd/server.go && \
  go build -tags 'netgo' -gcflags 'all=-l=4' -ldflags 'all=-s' -trimpath -o /tmp/kawipiko/bin/kawipiko-archiver ./cmd/archiver.go

FROM alpine AS site-builder

RUN mkdir /output
RUN echo "<html><body>This is kawipiko</body></html>" > ./output/index.html
COPY --from=kawipiko-builder /tmp/kawipiko/bin/kawipiko-archiver /usr/bin
RUN kawipiko-archiver --archive site.cdb --sources ./output --debug


# Switch to the scratch image
FROM scratch
EXPOSE 8080
# Copy the static binary
COPY --from=kawipiko-builder /tmp/kawipiko/bin/kawipiko-server /
# Copy the static website
COPY --from=site-builder /site.cdb .
# Run kawipiko
CMD ["/kawipiko-server", "--bind", "0.0.0.0:8080", "--archive", "/site.cdb", "--archive-mmap", "--report"]
