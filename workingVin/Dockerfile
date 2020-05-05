FROM golang:1.14-alpine AS build

RUN apk update \
    && apk add --virtual build-dependencies \
        build-base \
        gcc \
        wget \
        git \
        openssl-dev \
    && apk add bash pkgconfig

WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN ./configure && make && make install
WORKDIR /usr/local/lib
RUN cp librdkafka.a librdkafka-static.a
ADD . /go/src/github.com/kyungseopkim/goworkbed/workingVin
WORKDIR /go/src/github.com/kyungseopkim/goworkbed/workingVin
RUN go get
RUN go install -tags=static

RUN rm -rf /root/librdkafka

FROM golang:1.14-alpine AS final
COPY --from=build /go/bin/workingVin /go/bin/workingVin
ENTRYPOINT [ "/go/bin/workingVin" ]
