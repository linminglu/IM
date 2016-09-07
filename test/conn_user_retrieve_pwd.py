#-*- coding: utf8 -*-

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

# (errCode, strJson) = conn_test_com.TcpLogin(s, platform='a', phonenum="+8613812343333", passwd="25d55ad283aa400af464c76d713c07ad")
#
# print "login code:", errCode
# print "json:", strJson
#
# reto = json.loads(strJson)
# print reto
#
# sid = reto[u'sid']
# uid = reto[u'uid']

conn_test_com.TcpRetrievePwd(s, uid=100035, sid=0, phonenum='+8613812343333', password='25d55ad283aa400af464c76d713c07ad')

s.close()
time.sleep(4)

