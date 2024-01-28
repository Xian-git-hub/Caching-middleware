#!/bin/bash
# Start Redis server
redis-server /usr/local/etc/redis/redis.conf &


cp -r /usr/local/etc/caching-middleware/setting /app

cd /app
cache_middleware
