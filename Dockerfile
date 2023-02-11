FROM alpine

RUN apk add --no-cache git && \
    mkdir /workdir && \
    git config --global --add safe.directory /workdir

COPY git-chglog /usr/local/bin/git-chglog

WORKDIR /workdir
RUN chmod +x /usr/local/bin/git-chglog

ENTRYPOINT [ "/usr/local/bin/git-chglog" ]