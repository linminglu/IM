
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

uid = 100165

(errCode, strJson)  = conn_test_com.TcpLogin(uid,  "test2@qq.com", 'test', 'ahdjfhasjdfhaksdjfh324adsf8', s )
print "login code:", errCode
print "json:", strJson

if int(errCode) != 0 :
    print "login fail" 
    exit(0)
    
reto = json.loads(strJson)
print reto

sid = reto[u'sid']

#(errCode, strjson)  = conn_test_com.TcpSetUserInfo(100005, "Hello imsdk", sid, s )
#print "hello code:", errCode
#print "json:", strjson

req_json = '{"type":1}'
(errCode, strJson)  = conn_test_com.TcpSendReq2(50027, uid , req_json, sid, s)
print "create team code:", errCode
print "json:", strJson
