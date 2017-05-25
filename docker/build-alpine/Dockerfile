FROM golang:1.7.5-alpine

LABEL maintainer "cheminlin@cepave.com, robywei@cepave.com"

ENV \
  TZ="Asia/Taipei" \
  GET_NAME="github.com/Cepave/open-falcon-backend" \
  GIT_BRANCH="develop" \
  GIT_TAG="" \
  INTERVAL="180"
  # no_proxy="169.254/16, localhost, 127.0.0.1" \
  # AWS_KEY="" \
  # AWS_SECRET="" \
  # AWS_S3_BUCKET="" \
  # FTP_URL="" \
  # FTP_AUTH=""

RUN \
  apk update && apk add --no-cache bash ncurses make g++ wget tzdata git curl \
  && cp /usr/share/zoneinfo/"$TZ" /etc/localtime && echo "$TZ" > /etc/timezone \
  && cd /go/src && go get "$GET_NAME" \
  && cd /go/src/"$GET_NAME" && git checkout "$GIT_BRANCH" && git checkout "$GIT_TAG"

WORKDIR /go/src/"$GET_NAME"

## Start
ENTRYPOINT ["/bin/bash", "docker/build-alpine/scripts/build.sh"]
