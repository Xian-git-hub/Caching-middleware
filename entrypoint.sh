#!/bin/bash
# Start Redis server
redis-server /usr/local/etc/redis/redis.conf &

# 启动您的应用程序
cache_middleware
