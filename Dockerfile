FROM golang:1.21 as builder

WORKDIR /usr/src/app
COPY . .
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0
RUN go mod download && go mod verify
RUN go build -o cache_middleware .


FROM redis:6.2
COPY redis.conf /usr/local/etc/redis/redis.conf 

WORKDIR /app
VOLUME [ "/app" ]
WORKDIR /app
COPY --from=builder /usr/src/app/setting setting
COPY --from=builder /usr/src/app/cache_middleware /bin/cache_middleware
COPY --from=builder /usr/src/app/entrypoint.sh entrypoint.sh

EXPOSE 8080

RUN chmod +x entrypoint.sh
ENTRYPOINT [ "./entrypoint.sh" ]