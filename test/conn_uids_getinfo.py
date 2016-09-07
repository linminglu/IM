
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

#(errCode, strjson)  = conn_test_com.TcpSetUserInfo(100005, "Hello imsdk", sid, s )
#print "hello code:", errCode
#print "json:", strjson

req_json = '''{ "propertylist" : ["cid", "v"],
  "uidlist" : [
    16780303007937,
    16782433714609,
    18987547100514,
    18987513546082,
    18987479991650,
    18987463214434,
    18987446437217
    ]}'''
  
(errCode, strJson)  = conn_test_com.TcpSendReq2(10105, uid , req_json, sid, s)
print "get team info:", errCode
print "json:", strJson


eam info: 0
json: {"infolist":[{ "cid":"4423_ent.activity4215", "v":"1427025933", "uid":18987479991650},{ "cid":"5847_ent.activity4692", "v":"1427025933", "uid":18987513546082},{ "cid":"14319_ent.activity1584", "v":"1427025933", "uid":18987547100514},{ "cid":"test", "v":"1427028133", "uid":16780303007937},{ "cid":"guowx002", "v":"1427025933", "uid":16782433714609},{ "cid":"1024067_.mdby.motan2.961", "v":"1427025933", "uid":18987446437217},{ "cid":"5838_ent.activity4692", "v":"1427025933", "uid":18987463214434}]}


