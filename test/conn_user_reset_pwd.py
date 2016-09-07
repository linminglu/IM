#-*- coding: utf8 -*-

import socket
import conn_test_com
import time

#C->S:
#    sContent = dwUin+dwType+wLen+K(dwUin+strmd5(passwd)+ dwStatus + strReserved)
#S->C:
#    sContent = cResult + stOther

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

conn_test_com.TcpResetPwd(s, uid=100017, sid=10018, oldpwd='e10adc3949ba59abbe56e057f20f883e', newpwd='e10adc3949ba59abbe56e057f20f883a')

s.close()
time.sleep(4)

