FROM docker.io/alpine:latest

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY git-chglog /usr/local/bin/git-chglog

WORKDIR /workdir
RUN chmod +x /usr/local/bin/git-chglog

ENTRYPOINT [ "/usr/local/bin/git-chglog" ]
CMD [ "--help" ]
