FROM golang:alpine as builder
WORKDIR /home
ADD . .
RUN go build --mod=vendor -o axisblog main.go

FROM alpine as runner
ENV WORKDIR=/home
WORKDIR $WORKDIR

RUN apk add tzdata --no-cache \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
RUN mkdir -p $WORKDIR/static #创建静态文件目录
COPY --from=builder $WORKDIR/axisblog $WORKDIR
COPY --from=builder $WORKDIR/config.toml $WORKDIR
COPY --from=builder $WORKDIR/sensitive-words.txt $WORKDIR
CMD ["./axisblog"]