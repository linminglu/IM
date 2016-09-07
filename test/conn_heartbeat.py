#-*- coding: utf8 -*-

import socket
import conn_test_com

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))
data = conn_test_com.TcpHello(167772160001, 10001, s)
print data
