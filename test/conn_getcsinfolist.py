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

for ip in ["122.13.81.199","14.29.84.61"]:    
    print ip, "==========="
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((ip, conn_test_com.CONN_PORT))
    print (ip, conn_test_com.CONN_PORT)
    (errCode, strJson)  = conn_test_com.TcpLogin(100149, 'IMUser167', '96e79218965eb72c92a549dd5a330112', '00b6413a92d4c1c84ad99e0a', s )
    print "login code:", errCode
    print "json:", strJson
    
    
    if errCode != 0 :
        sys.exist(0)
    
    
    reto = json.loads(strJson)
    print reto
    
    
    sid = reto[u'sid']
    uid = reto[u'uid']
    
#(errCode, strjson)  = conn_test_com.TcpSetUserInfo(100005, "Hello imsdk", sid, s )
#print "hello code:", errCode
#print "json:", strjson

    req_json = '{"appkey":"00b6413a92d4c1c84ad99e0a"}'
    (errCode, strJson)  = conn_test_com.TcpSendReq2(11003, uid , req_json, sid, s)
    print "set team info:", errCode
    print "json:", strJson

