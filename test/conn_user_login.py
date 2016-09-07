#-*- coding: utf8 -*-

import socket
import conn_test_com
import time
import json
import sys
import StringIO
import struct

#C->S:
#    sContent = dwUin+dwType+wLen+K(dwUin+strmd5(passwd)+ dwStatus + strReserved)
#S->C:
#    sContent = cResult + stOther

# for ip in ["122.13.81.199", "14.29.84.61"]:
# for ip in ["120.25.239.220"]:
# for ip in ["127.0.0.1"]:
# print ip, "==========="

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

print (conn_test_com.CONN_IP, conn_test_com.CONN_PORT)

(errCode, strJson) = conn_test_com.TcpLogin2(s, platform='a', uid=3000071, passwd="1xxswebd")

print "login code:", errCode
print "json:", strJson

if errCode != 0:
    sys.exit(0)

reto = json.loads(strJson)
print reto

sid = reto[u'sid']
uid = reto[u'uid']

time.sleep(3)

# conn_test_com.TcpHello(uid, sid, s)

# print "errCode=%s strjson=%s" % (errCode, strjson)

# recv chat msg
msgids=[]

    # while 1:
    #     s.settimeout(3)
    #     havemsg = 0
    #     try:
    #         print "recv dbmsg"
    #         data = s.recv(10240)
    #         print "read buf size:",  len(data)
    #     except socket.error, msg:
    #         print msg
    #         break
    #     packlen = 0
    #     while  len(data) > packlen :
    #         tempstr = StringIO.StringIO(data)
    #         tempstr.seek(0, 0)
    #
    #         pkghead =  tempstr.read(24)
    #
    #         fmt = "!HHHHIQI"
    #         (Len, Cmd, version, Seq , errCode, uid, compress)= struct.unpack(fmt, pkghead)
    #         print (Len , errCode, Cmd , Seq)
    #         packlen = packlen + Len
    #
    #         jsonlen =  Len - 24
    #
    #         fmt = "%ds" % (jsonlen)
    #         jsonbody  =  tempstr.read(jsonlen)
    #         (strjson,) = struct.unpack(fmt, jsonbody)
    #         print strjson
    #         reto = json.loads(strjson)
    #         if Cmd == 30101:
    #             msgid = reto["text_chat_msg"]['MsgId']
    #             errcode , retjson = conn_test_com.TcpMsgRecved( 100150, sid,  int(msgid), s)
    #             print "recv:"  ,errcode ,  retjson
    #         elif Cmd == 35100:
    #             print "recvesysmg", strjson
    #         elif Cmd == 35101:
    #             print "recveteammg", strjson
    #             msgid = reto['msgid']
    #             print "msgid:", msgid
    #             msgids.append(msgid)
    #
    # if len(msgids) > 0 :
    #     for msgid in msgids :
    #         errcode , retjson = conn_test_com.TcpUserMsgReceived(uid, sid, s, msgid, 0)
    #         print "recv:"  ,errcode ,  retjson

while True:
    time.sleep(10)
    s.settimeout(2)
    havemsg = 0
    try:
        conn_test_com.TcpHello(uid, sid, s)
        print "recv msg"
        data = s.recv(10240)
        print "read buf size:", len(data)
    except socket.error, msg:
        print msg
        # break

    # packlen = 0
    # while len(data) > packlen:
    #     tempstr = StringIO.StringIO(data)
    #     tempstr.seek(0, 0)
    #
    #     pkghead = tempstr.read(24)
    #
    #     fmt = "!HHHHIQI"
    #     (Len, Cmd, version, Seq, errCode, uid, compress) = struct.unpack(fmt, pkghead)
    #
    #     print "Len=%d,errCode=%d,Cmd=%d,Seq=%s" % (Len, errCode, Cmd, Seq)
    #
    #     packlen = packlen + Len
    #
    #     jsonlen = Len - 24
    #
    #     fmt = "%ds" % (jsonlen)
    #     jsonbody = tempstr.read(jsonlen)
    #     (strjson, ) = struct.unpack(fmt, jsonbody)
    #
    #     print "strjson=", strjson
    #
    #     reto = json.loads(strjson)
    #     if Cmd == 30101:
    #         msgid = reto["text_chat_msg"]['MsgId']
    #         errcode, retjson = conn_test_com.TcpMsgRecved(100150, sid, int(msgid), s)
    #         print "recv:", errcode,  retjson
    #     elif Cmd == 35100:
    #         print "recvesysmg", strjson
    #     elif Cmd == 35101:
    #         print "recveteammg", strjson
    #         msgid = reto['msgid']
    #         print "msgid:", msgid
    #         msgids.append(msgid)

if len(msgids) > 0:
    for msgid in msgids:
        errcode, retjson = conn_test_com.TcpUserMsgReceived(uid, sid, s, msgid, 0)
        print "recv:", errcode, retjson
