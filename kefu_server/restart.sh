#! /bin/sh
#
# restart.sh
# Copyright (C) 2014 nh <nh@gw.kxc.imsdk.im>
#
# Distributed under terms of the MIT license.
#


killall -9 restful_svr
sleep 2
./restful_svr svr.conf &
