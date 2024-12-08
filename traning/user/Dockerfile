FROM alpine

# RUN mkdir /go && cd /go \
#   && wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz \
#   && tar -C /usr/local -zxf go1.23.1.linux-amd64.tar.gz \
#   && rm -rf /go/go1.23.1.linux-amd64.tar.gz \
#   && mkdir /lib64 \
#   && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# ENV GOPATH=/go
# ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

# CMD ["ping", "www.baidu.com"]


# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/user-api ./api/user.go


# docker build -t user-api -f ./Dockerfile .

RUN mkdir /user

COPY bin/user-api /user

COPY api/etc/user.yaml /user

RUN chmod +x /user/user-api

WORKDIR /user

ENTRYPOINT ["./user-api", "-f", "user.yaml"]


# docker run --name user-api -p 8888:8888 --net user --link user-etcd -d user-api