#!/bin/bash
# Start Redis server
redis-server /usr/local/etc/redis/redis.conf &


if [ ! -d "/app/setting" ]; then
  cp -r /usr/local/etc/caching-middleware/setting /app
fi

cache_middleware
