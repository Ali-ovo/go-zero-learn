FROM alpine

RUN mkdir /go && cd /go \
  && wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz \
  && tar -C /usr/local -zxf go1.23.1.linux-amd64.tar.gz \
  && rm -rf /go/go1.23.1.linux-amd64.tar.gz \
  && mkdir /lib64 \
  && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENV GOPATH=/go
ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

CMD ["ping", "www.baidu.com"]