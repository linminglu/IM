
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

(errCode, strJson) = conn_test_com.TcpLogin(s, platform='a', phonenum=conn_test_com.TEST_PHONE_NUM, cid='test', passwd=conn_test_com.TEST_PWD)

print "login code:", errCode
print "json:", strJson

reto = json.loads(strJson)
print reto

sid = reto[u'sid']
uid = reto[u'uid']

req_json = '{"msgid":1410771020985758, "uid":100079}'
(errCode, strJson) = conn_test_com.TcpSendReq2(30201, 100079, req_json, sid, s)

print "get team info:", errCode
print "json:", strJson

sid = reto[u'sid']
msgid = 1410774092487913

errcode, retjson = conn_test_com.TcpUserMsgReceived(100032, sid, s, msgid, 100033)

print "recv:", errcode, retjson

