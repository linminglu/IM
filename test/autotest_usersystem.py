#-*- coding: utf8 -*-

import socket
import conn_test_com
import string
import json
import sys
import comm

def HeartBeat():
    print "------------------------Test Heart Beat--------------------------------"

    errCode = conn_test_com.TcpHello(comm.gl_uid, comm.gl_sid, s)

    print errCode

def Register():
    print "------------------------Test Register--------------------------------"

    print "Input platform: (i for ios, a for android, other for default):"
    line = sys.stdin.readline()
    comm.gl_platform = line[:-1]

    print "Input appkey:"
    line = sys.stdin.readline()
    comm.gl_appkey = line[:-1]

    print "Input cid:"
    line = sys.stdin.readline()
    comm.gl_cid = line[:-1]

    print "Input did:"
    line = sys.stdin.readline()
    did = line[:-1]

    print "Input password:"
    line = sys.stdin.readline()
    comm.gl_password = line[:-1]

    (errCode, strJson)  = conn_test_com.TcpReg(comm.gl_appkey, comm.gl_cid, comm.gl_password, s)
    print "reg code:", errCode
    print "json:", strJson
    reto = json.loads(strJson)
    print reto

    comm.gl_uid = reto[u'uid']
    comm.gl_password = reto[u'password']
    comm.gl_sid = reto[u'sid']

	
def Login():
    print "------------------------Test Login--------------------------------"

    print "Input platform: (i for ios, a for android, other for default):"
    line = sys.stdin.readline()
    comm.gl_platform = line[:-1]

    if comm.gl_platform != 'i' and comm.gl_platform != 'a':
        print "login use default user: ", comm.gl_uid

    else:
        print "Input uid:"
        line = sys.stdin.readline()
        comm.gl_uid = string.atoi(line[:-1])

        print "Input cid:"
        line = sys.stdin.readline()
        comm.gl_cid = line[:-1]

        print "Input password:"
        line = sys.stdin.readline()
        comm.gl_password = line[:-1]

        print "Input appkey:"
        line = sys.stdin.readline()
        comm.gl_appkey = line[:-1]

    (errCode, strJson)  = conn_test_com.TcpLogin(comm.gl_uid, comm.gl_cid, comm.gl_password, comm.gl_appkey, s )
    print "login code:", errCode
    print "json:", strJson
    reto = json.loads(strJson)
    print reto

    comm.gl_uid = reto[u'uid']
    comm.gl_sid = reto[u'sid']

def SetDeviceToken():
	print "------------------Test Set Device Token--------------------------------"

	print "Input device token:"
	line = sys.stdin.readline()
	token = line[:-1]

	(errCode, strJson) = conn_test_com.TcpSetDevToken(comm.gl_uid, token, comm.gl_sid, s)
	print "set device token code:", errCode


def SetUserInfo():
    print "------------------------Test Set User Info-----------------------------"

    print "Input did:"
    line = sys.stdin.readline()
    did = line[:-1]

    print "Input baseinfo:"
    line = sys.stdin.readline()
    baseinfo = line[:-1]

    print "Input exinfo:"
    line = sys.stdin.readline()
    exinfo = line[:-1]

    (errCode, strJson)  = conn_test_com.TcpSetUserInfo(comm.gl_uid, comm.gl_sid, s, did, baseinfo, exinfo)
    print "set user info code:", errCode
    print "json:", strJson


def GetUserInfo():
    print "------------------------Test Get User Info-----------------------------"

    uidlist = []
    print "Please input uid: (non-number for end)"
    while True:
        line = sys.stdin.readline()
        if line[:-1].isdigit():
            uid = string.atoi(line[:-1])
            uidlist.append(uid)
        else:
            break

    propertylist = []
    print "Please input property: (null for end)"
    while True:
        line = sys.stdin.readline()
        if line[:-1] != "":
            propertylist.append(line[:-1])
        else:
            break

    (errCode, strJson)  = conn_test_com.TcpGetUserInfo(uidlist, propertylist, comm.gl_sid, s)
    print "get user info code:", errCode
    print "json:", strJson

def GetUid():
	print "------------------------Test Get Uid-----------------------------"

	cidlist = []
	print "Please input cid: (null for end)"
	while True:		
		line = sys.stdin.readline()
		if line[:-1] != "":			
			cidlist.append(line[:-1])
		else:
			break

	(errCode, strJson)  = conn_test_com.TcpGetUid(comm.gl_uid, comm.gl_appkey, cidlist, comm.gl_sid, s)
	print "get user info code:", errCode
	print "json:", strJson

def SaveLocation():
	print "------------------------Test Save Location-----------------------------"

	print "Input XPos:"
	line = sys.stdin.readline()
	xpos = string.atof(line[:-1])

	print "Input YPos:"
	line = sys.stdin.readline()
	ypos = string.atof(line[:-1])

	(errCode, strJson)  = conn_test_com.TcpUserLocation(comm.gl_uid, xpos, ypos, comm.gl_sid, s, 0)
	print "save user info code:", errCode
	print "json:", strJson

def GetLocation():
	print "------------------------Test Get Location-----------------------------"

	print "Input XPos:"
	line = sys.stdin.readline()
	xpos = string.atof(line[:-1])

	print "Input YPos:"
	line = sys.stdin.readline()
	ypos = string.atof(line[:-1])

	print "Input level:(greater than 0)"
	line = sys.stdin.readline()
	level = string.atoi(line[:-1])

	(errCode, strJson)  = conn_test_com.TcpUserLocation(comm.gl_uid, xpos, ypos, comm.gl_sid, s, level)
	print "get user info code:", errCode
	print "json:", strJson

	
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((conn_test_com.CONN_IP, conn_test_com.CONN_PORT))

print "------------------------Test Start--------------------------------"
print "1.Heart Beat"
print "2.Register"
print "3.Login"
print "4.Set Device Token"
print "5.Set User Info"
print "6.Get User Info"
print "7.Get Uid"
print "11.Save Location"
print "12.Get Location"
print "q.Quit Test"

while True:
	print "Please input test module number: (q for quit)"

	line = sys.stdin.readline()
	if 'q' == line[:-1]:
		break

	moduleNum = string.atoi(line[:-1])
	{
		1: lambda: HeartBeat(),
		2: lambda: Register(),
		3: lambda: Login(),
		4: lambda: SetDeviceToken(),
		5: lambda: SetUserInfo(),
		6: lambda: GetUserInfo(),
		7: lambda: GetUid(),

		11: lambda: SaveLocation(),
		12: lambda: GetLocation(),
	}[moduleNum]()

	print comm.gl_uid, comm.gl_sid

print "------------------------Test Finished-----------------------------"