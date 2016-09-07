#!/bin/sh

nsqlookupd="nsqlookupd"
nsqd="nsqd"
nsqadmin="nsqadmin"
log_file="nsqd.log"
DEV_MODE="debug"

start() {
	mkdir ".nsqd_data" 2>>/dev/null 1>>/dev/null
	if [ "$DEV_MODE" = "debug" ]; then
		$nsqlookupd &  2>>/dev/null 1>>/dev/null                                          
		$nsqadmin --lookupd-http-address=127.0.0.1:4161 &  2>>/dev/null 1>>/dev/null  
		$nsqd -data-path ".nsqd_data" --lookupd-tcp-address=127.0.0.1:4160 & 2>>/dev/null 1>>/dev/null        
	else
		$nsqd -data-path ".nsqd_data" &        
	fi
}

stop() {
	killall "$nsqd" 2>/dev/null
	if [ "$DEV_MODE" = "debug" ]; then
		killall "$nsqlookupd" 2>/dev/null
		killall "$nsqadmin" 2>/dev/null
	fi
}

restart() {
	stop
	start
}
status_p() {
	status "$nsqlookupd" 
	status "$nsqd"       
	status "$nsqadmin"   
}

case "$1" in 
	"start")
		start
		;;
	"stop")
		stop
		;;
	"restart")
		restart
		;;
	*)
		echo "Usage: $0 {start|stop|restart}"
		exit 2
		;;
esac




