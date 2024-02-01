FROM golang:1.21 as builder

WORKDIR /usr/src/app
COPY . .
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0
RUN go mod download && go mod verify
RUN go build -o cache_middleware .


FROM redis:6.2
WORKDIR /app
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

COPY redis.conf /usr/local/etc/redis/redis.conf 
COPY --from=builder /usr/src/app/cache_middleware /bin/cache_middleware

COPY --from=builder /usr/src/app/setting /usr/local/etc/caching-middleware/setting
COPY --from=builder /usr/src/app/entrypoint.sh /usr/local/etc/caching-middleware/entrypoint.sh

EXPOSE 8080

RUN chmod +x /usr/local/etc/caching-middleware/entrypoint.sh
ENTRYPOINT [ "/usr/local/etc/caching-middleware/entrypoint.sh" ]
