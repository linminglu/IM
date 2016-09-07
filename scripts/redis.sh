#!/bin/sh

# chkconfig: 2345 10 90
# description: Start and Stop redis

REDISPORT=6379                                  
EXC=/usr/local/redis/bin/redis-server
REDIS_CLI=/usr/local/redis/bin/redis-cli

PIDFILE=/var/tmp/redis.pid
CONF="/usr/local/redis/etc/redis.conf"              

case "$1" in
        start)
                echo "Starting Redis server..."
                $EXC $CONF
				;;
        stop)
                echo "Stoping Redis server..."
                $REDIS_CLI -p $REDISPORT SHUTDOWN
                ;;
        restart|force-reload)
                ${0} stop
                ${0} start
                ;;
        *)
                echo "Usage: /etc/init.d/redis {start|stop|restart|force-reload}" >&2
                exit 1
esac

exit 0
