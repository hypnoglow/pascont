FROM golang:1.8-alpine
LABEL maintainer "Igor Zibarev <zibarev.i@gmail.com>"

RUN set -ex \
    && apk add --no-cache --virtual .build-deps git

WORKDIR /go/src/github.com/hypnoglow/pascont
COPY . .

RUN go-wrapper download
RUN go-wrapper install

RUN set -ex \
    && apk del .build-deps

RUN go test $(go list ./... | grep -v "/vendor/")

CMD ["go-wrapper", "run"]