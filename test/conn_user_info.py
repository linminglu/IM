#-*- coding: utf8 -*-

import socket
import conn_test_com
import json

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

(errCode, strJson) = conn_test_com.default_login(s)

print "login code:", errCode
print "json:", strJson

reto = json.loads(strJson)
print reto

sid = reto[u'sid']
uid = reto[u'uid']

req_json = '''{
    "uidlist": [100040,100041,100042,100043],
    "vlist": [1, 0, 2],
    "propertylist": ["phonenum"]
}'''

for i in range(0, 1):
    (errCode, strJson) = conn_test_com.TcpSendReq2(s, cmd=10105, req_json=req_json, uid=uid, sid=sid)
    print "TcpSendReq2 code:", errCode
    print "json:", strJson
    print i, "---------------------"
