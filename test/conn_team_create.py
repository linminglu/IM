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

sid = reto[u'sid']
uid = reto[u'uid']

create_team_json = '{"teamname":"attention2"}'
(errCode, strJson) = conn_test_com.TcpSendReq2(s, 50001, create_team_json, uid, sid)
print "create team code:", errCode
print "json:", strJson
