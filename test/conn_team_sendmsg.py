#-*- coding: utf8 -*-

import socket
import conn_test_com
import json

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

(errCode, strJson) = conn_test_com.TcpLogin(s, platform='a', phonenum=conn_test_com.TEST_PHONE_NUM, passwd=conn_test_com.TEST_PWD)
print "login code:", errCode
print "json:", strJson

if int(errCode) != 0:
    print "login fail"
    exit(0)

reto = json.loads(strJson)
print reto

uid = reto[u'uid']
sid = reto[u'sid']

req_json = '{"teamid":100028001, "uid":100028}'
(errCode, strJson) = conn_test_com.TcpSendReq2(s, 50023, req_json, uid, sid)
print "set team info:", errCode
print "json:", strJson

