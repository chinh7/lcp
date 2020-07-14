FROM golang:1.14-alpine

RUN echo "@testing http://nl.alpinelinux.org/alpine/latest-stable/main" >>/etc/apk/repositories
RUN apk add --upda`te --no-cache build-base linux-headers git cmake bash #wget mercurial g++ autoconf libgflags-dev cmake bash
RUN apk add --update --no-cache zlib zlib-dev bzip2 bzip2-dev snappy snappy-dev lz4 lz4-dev zstd@testing zstd-dev@testing

RUN cd /tmp && \
  git clone https://github.com/gflags/gflags.git && \
  cd gflags && \
  mkdir build && \
  cd build && \
  cmake -DBUILD_SHARED_LIBS=1 -DGFLAGS_INSTALL_SHARED_LIBS=1 .. && \
  make install && \
  cd /tmp && \
  rm -R /tmp/gflags/

RUN cd /tmp && \
  git clone https://github.com/facebook/rocksdb.git && \
  cd rocksdb && \
  git checkout v6.2.2 && \
  make shared_lib && \
  mkdir -p /usr/local/rocksdb/lib && \
  mkdir /usr/local/rocksdb/include && \
  cp librocksdb.so* /usr/local/rocksdb/lib && \
  cp /usr/local/rocksdb/lib/librocksdb.so* /usr/lib/ && \
  cp -r include /usr/local/rocksdb/ && \
  cp -r include/* /usr/include/ && \
  rm -R /tmp/rocksdb/
