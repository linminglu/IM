
#-*- coding: utf8 -*-

import struct
import StringIO
import socket
import conn_test_com
import time
import json
import sys

#C->S:
#    sContent = dwUin+dwType+wLen+K(dwUin+strmd5(passwd)+ dwStatus + strReserved)
#S->C:
#    sContent = cResult + stOther


# for ip in ["122.13.81.199"]:#"14.29.84.61"]:
print conn_test_com.CONN_IP, "==========="

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

print (conn_test_com.CONN_IP, conn_test_com.CONN_PORT)

# (errCode, strJson) = conn_test_com.TcpLogin(s, platform='a', phonenum=conn_test_com.TEST_PHONE_NUM, passwd=conn_test_com.TEST_PWD)

(errCode, strJson) = conn_test_com.TcpLogin(s, platform='a', phonenum="+8613812341111", passwd=conn_test_com.TEST_PWD)
# (errCode, strJson) = conn_test_com.TcpLogin(s, platform='a', phonenum="+8618816822806", passwd=conn_test_com.TEST_PWD)

print "login code:", errCode
print "json:", strJson

if errCode != 0:
    sys.exist(0)

reto = json.loads(strJson)
print reto

sid = reto[u'sid']
uid = reto[u'uid']

for i in range(0, 1):
    time.sleep(1)

    msg = "你好"
    jsonStr = '{"msgtype":0,"touid":100042,"msgcontent": "%s","apnstext":"gwx 给您发送了一条消息"}' % (msg)

    conn_test_com.OnlyTcpSendReq2(s, cmd=30101, req_json=jsonStr, uid=uid, sid=sid)

    s.settimeout(5)
    isok = 0
    try:
        print "recv"
        data = s.recv(10240)
        isok = 11
    except socket.error, msg:
        print msg

    if isok > 0:
        tempstr = StringIO.StringIO(data)
        tempstr.seek(0, 0)

        pkghead = tempstr.read(len(data))

        fmt = "!HHHHIII%ds" % (len(data)-20)
        (Len,  Cmd, version, Seq, errCode, uid, Flag, strJson) = struct.unpack(fmt, pkghead)
        print (errCode, Cmd, Seq, strJson)

s.close()
