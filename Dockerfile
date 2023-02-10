##################
########## Builder
##################
FROM golang:1.19-alpine as builder

RUN apk add --no-cache tzdata make git

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

##################
########## PROD
##################
FROM alpine

RUN apk add --no-cache git && \
    mkdir /workdir

COPY --from=builder /build/git-chglog /usr/local/bin/git-chglog

WORKDIR /workdir
RUN chmod +x /usr/local/bin/git-chglog

ENTRYPOINT [ "/usr/local/bin/git-chglog" ]