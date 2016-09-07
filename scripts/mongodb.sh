#!/bin/sh

# chkconfig: 2345 10 90
# description: Start and Stop nsqd


PROG=mongod

PIDFILE="/var/run/${PROG}.pid"

case "$1" in
        start)
#                if [ -f $PIDFILE ]
#                then
#                        echo "$PIDFILE exists, ${PROG} is already running or crashed."
#                else
                        echo "Starting ${PROG} server..."
                        $PROC &
#                fi
#                if [ "$?"="0" ]
#                then
#                        touch $PIDFILE
#                        echo "${PROG} is running..."
#                fi
                ;;
        stop)
#                if [ ! -f $PIDFILE ]
#                then
#                        echo "$PIDFILE exists, ${PROG} is not running."
#                else
#                        echo "$PROG Stopping..."
                        killall $PROG
#                        while [ -x $PIDFILE ]
#                        do
#                                echo "Waiting for $PROG to shutdown..."
#                                sleep 1
#                        done
#                        #rm -f $PIDFILE
#                        echo "${PROG} stopped"
#                fi
                ;;

        restart|force-reload)
                ${0} stop
                ${0} start
                ;;
        *)
                echo "Usage: /etc/init.d/$PROG {start|stop|restart|force-reload}" >&2
                exit 1
esac

