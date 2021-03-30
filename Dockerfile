FROM alpine:3.13.3

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY git-chglog /usr/local/bin/git-chglog

WORKDIR /workdir
RUN chmod +x /usr/local/bin/git-chglog

ENTRYPOINT [ "/usr/local/bin/git-chglog" ]