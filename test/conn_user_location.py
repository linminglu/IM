
#-*- coding: utf8 -*-

import struct
import StringIO
import socket
import conn_test_com
import time
import json
#C->S:
#    sContent = dwUin+dwType+wLen+K(dwUin+strmd5(passwd)+ dwStatus + strReserved)
#S->C:
#    sContent = cResult + stOther


s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

(errCode, strJson)  = conn_test_com.TcpLogin(16777266331841, 'nnhh', '453e41d218e071ccfb2d1c99ce23906a', '16777266331841', s )
print "login code:", errCode
print "json:", strJson

reto = json.loads(strJson)
print reto

sid = reto[u'sid']

(errCode, strJson)  = conn_test_com.TcpUserLocation(16777266331841, 116.207629, 41.058359, sid, s, 1 , page=0)
print "get user info code:", errCode
print "json:", strJson
if errCode == 0 :
    ret = json.loads(strJson)
    for it in ret['userlocations']:
            print it
