FROM golang:1.14.0-alpine AS build

RUN apk update \
    && apk add --virtual build-dependencies \
        build-base \
        gcc \
        wget \
        git \
    && apk add bash pkgconfig

WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN ./configure && make && make install

ADD . /go/src/github.com/kyungseopkim/goworkbed/kafka
WORKDIR /go/src/github.com/kyungseopkim/goworkbed/kafka
RUN go get
RUN go install -tags=static

RUN rm -rf /root/librdkafka

FROM golang:1.14-alpine AS final
COPY --from=build /go/bin/kafka /go/bin
ENTRYPOINT [ "/go/bin/kafka" ]
