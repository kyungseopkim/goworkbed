FROM golang:1.14.2-alpine AS build
RUN apk update && apk add --no-cache  \
        bash              \
        build-base        \
        coreutils         \
        gcc               \
        git               \
        make              \
        musl-dev          \
        openssl-dev       \
        openssl           \
        libsasl           \
        libgss-dev        \
        rpm               \
        lz4-dev           \
        zlib-dev          \
        ca-certificates   \
        wget              \
        curl          &&  \
    cd $(mktemp -d) && \
    wget -nv -O cyrus-sasl-2.1.27.tar.gz https://github.com/cyrusimap/cyrus-sasl/releases/download/cyrus-sasl-2.1.27/cyrus-sasl-2.1.27.tar.gz && \
    tar -xz --strip-components=1 -f cyrus-sasl-2.1.27.tar.gz && \
    rm -f cyrus-sasl-2.1.27.tar.gz && \
    ./configure --prefix=/usr --disable-sample --disable-obsolete_cram_attr --disable-obsolete_digest_attr --enable-static --disable-shared \
        --disable-checkapop --disable-cram --disable-digest --enable-scram --disable-otp --disable-gssapi --with-dblib=none --with-pic && \
    make && \
    make install

WORKDIR /root
RUN cd $(mktemp -d) && \
    wget -nv -O v1.4.0.tar.gz https://github.com/edenhill/librdkafka/archive/v1.4.0.tar.gz && \
    tar -xz --strip-components=1 -f v1.4.0.tar.gz && \
    rm -f v1.4.0.tar.gz && \
    ./configure --prefix=/usr --enable-sasl && \
    make -j && \
    make install

RUN  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ADD . /go/src/github.com/kyungseopkim/goworkbed/seqsinker
WORKDIR /go/src/github.com/kyungseopkim/goworkbed/seqsinker
RUN dep init
RUN go build -tags dynamic -o seqsinker && cp seqsinker /go/bin

FROM golang:1.14.2-alpine AS final
COPY --from=build /usr/lib/librdkafka* /usr/lib/
COPY --from=build /usr/lib/liblz4* /usr/lib/
COPY --from=build /go/bin/seqsinker /go/bin/

ENTRYPOINT [ "/go/bin/seqsinker" ]
