#-*- coding: utf8 -*-

import socket
import conn_test_com
import time

#C->S:
#    sContent = dwUin+dwType+wLen+K(dwUin+strmd5(passwd)+ dwStatus + strReserved)
#S->C:
#    sContent = cResult + stOther

# for i in range(0, 1):
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

(errCode, strJson) = conn_test_com.TcpRegister(s, platform='a', phonenum=conn_test_com.TEST_PHONE_NUM, passwd=conn_test_com.TEST_PWD)

print "login code:", errCode
print "json:", strJson

s.close()
time.sleep(4)

