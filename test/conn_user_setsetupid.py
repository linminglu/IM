
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

uid = 16780303007937
(errCode, strJson)  = conn_test_com.TcpLogin(uid, 'test', '098f6bcd4621d373cade4e832627b4f6', '00b6413a92d4c1c84ad99e0a', s)
print "login code:", errCode
print "json:", strJson


reto = json.loads(strJson)
print reto


sid = reto[u'sid']

for i in range(0, 1):
    time.sleep(1)
    msg = "msg%05d" % i
    req_json = '{"setupid":6113685117481206520}'
    #req_json = '{"msgtype":0, "touid":16777266331841 ,"msgcontent": "%s"}' % (msg)
    conn_test_com.OnlyTcpSendReq2(11002, uid , req_json, sid, s)

    s.settimeout(5)
    isok = 0
    try:
        print "recv"
        data = s.recv(10240)
        isok = 11
    except socket.error, msg:
        print msg
    
    
    if isok > 0 :
        tempstr = StringIO.StringIO(data)
        tempstr.seek(0, 0)
        
        pkghead =  tempstr.read(len(data))
        
        fmt = "!HHHHIII%ds" % (len(data)-20)
        (Len,  Cmd, version , Seq , errCode, uid , Flag , strJson)= struct.unpack(fmt, pkghead)
        print (errCode, Cmd , Seq, strJson)
        
s.close()
