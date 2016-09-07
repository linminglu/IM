
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


(errCode, strJson)  = conn_test_com.TcpLogin(100149, 'IMUser167', '96e79218965eb72c92a549dd5a330112', '00b6413a92d4c1c84ad99e0a', s )
print "login code:", errCode
print "json:", strJson

reto = json.loads(strJson)
print reto

sid = reto[u'sid']

#(errCode, strjson)  = conn_test_com.TcpSetUserInfo(100033, sid, s, "", "Hello imsdk", "By aimen xinxi keji youxian gongsi")
#print "hello code:", errCode
#print "json:", strjson

#(errCode, strjson)  = conn_test_com.TcpGetUserInfo(100033, sid, s)
#print "get user info code:", errCode
#print "json:", strjson
for i in range(1, 20):
    (errCode, strJson)  = conn_test_com.TcpSetDevToken(100149, '11111111111111111111111111221111', sid, s)
    print "get user info code:", errCode
    print "json:", strJson
