#!/bin/sh

case $1 in
	"start")
		sh mongodb.sh start
		sh nsqd.sh start
		sh redis.sh start
		;;
	"stop")
		sh mongodb.sh stop
		sh nsqd.sh stop
		sh redis.sh stop
		;;
	"restart")
		sh mongodb.sh restart 
		sh nsqd.sh restart
		sh redis.sh restart
		;;
esac
