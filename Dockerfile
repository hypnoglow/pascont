FROM golang:1.8-alpine
LABEL maintainer "Igor Zibarev <zibarev.i@gmail.com>"

RUN set -ex \
    && apk add --no-cache --virtual .build-deps git \
    && go get -u -v github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/hypnoglow/pascont

COPY . .

RUN dep ensure -v \
    && go-wrapper install

RUN set -ex \
    && apk del .build-deps

CMD ["go-wrapper", "run"]