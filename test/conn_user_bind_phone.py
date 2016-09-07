#-*- coding: utf8 -*-

import socket
import conn_test_com
import time

#C->S:
#    sContent = dwUin+dwType+wLen+K(dwUin+strmd5(passwd)+ dwStatus + strReserved)
#S->C:
#    sContent = cResult + stOther

for i in range(0, 1):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

    (errCode, strJson) = conn_test_com.TcpBindPhone(s, uid=conn_test_com.TEST_UID, sid=conn_test_com.TEST_SID, phonenum='+8613812343333')

    print "login code:", errCode
    print "json:", strJson

    s.close()
    time.sleep(4)

