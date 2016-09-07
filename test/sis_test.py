#!/home/push/python/bin/python

import socket
import sys
import os
import struct

a = '''
typedef struct  _UdpSisInfo
{
    char len[2];
    char head[2];
    char net_type[30];
    UInt tel_opera;
    UInt uid;
    char senderid[50];
   char sdk_ver[10];
    int  test_mode;
#    char res[22];
#}UdpSisInfo;'''


head = " E"
net_type = "NH-test                       "
tel_opera = 60001
uid = 65000111
senderid = 'qswddddddsddddsdwwssdsd     '
sdk_ver = '0.2.0     '
test_mode = 1
res = '                      '
pkglen = 2 + len(head) + len(net_type) + 4*2 + len(senderid) +  len(sdk_ver) + 4 + len(res)
print "pkglen", pkglen
pkg = struct.pack("!H2s30sII50s10sI22s", pkglen , head, net_type, tel_opera, uid, senderid , sdk_ver, test_mode, res)

address = ('120.25.239.220', 19000)
# address = ('127.0.0.1', 19000)

s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
# s.sendto(pkg, address)
s.sendto("", address)

s.settimeout(3)
try:
    data, addr = s.recvfrom(2048)
    print data
except :
    print  "timeout"   
 
    