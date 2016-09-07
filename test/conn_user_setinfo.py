
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


uid = 100005
(errCode, strJson)  = conn_test_com.TcpLogin(100149, 'IMUser167', '96e79218965eb72c92a549dd5a330112', '00b6413a92d4c1c84ad99e0a', s )
print "login code:", errCode
print "json:", strJson

if int(errCode) != 0 :
    print "login fail" 
    exit(0)
    
reto = json.loads(strJson)
print reto

sid = reto[u'sid']
uid = reto[u'uid']


req_json = '''{
    "exinfo": "exinfo set test",
    "baseinfo": "niehao test baseinfo"
    
}'''

(errCode, strJson)  = conn_test_com.TcpSendReq2(10104, uid , req_json, sid, s)
print "create team code:", errCode
print "json:", strJson
