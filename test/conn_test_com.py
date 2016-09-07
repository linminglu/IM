# -*- coding:utf-8 -*-

import struct
import socket
import StringIO

#推送包
CMD_PUSH = 5000

#三大系统
CMD_USER_SYSTEM = 10000
CMD_SOCIAL_SYSTEM =20000
CMD_IM_SYSTEM = 30000
CMD_APP_SYSTEM = 35000
#用户系统
#请求/应答包：
CMD_REGISTER =(CMD_USER_SYSTEM + 1)
CMD_LOGIN=(CMD_USER_SYSTEM + 2)
CMD_MODIFY_PASSWORD=(CMD_USER_SYSTEM + 3)
#define CMD_SET_PHOTO=(CMD_USER_SYSTEM + 4)
CMD_SET_MY_USER_INFO=(CMD_USER_SYSTEM + 5)
CMD_SET_LABS=(CMD_USER_SYSTEM + 12)

#社交系统
#请求/应答包
CMD_QUERY_USER_BASE_INFO = (CMD_SOCIAL_SYSTEM + 1)
CMD_QUERY_GROUP_USER_BASE_INFO = (CMD_SOCIAL_SYSTEM + 2)
CMD_QUERY_USER_INFO = (CMD_SOCIAL_SYSTEM + 3)
CMD_QUERY_USER_EXT_INFO = (CMD_SOCIAL_SYSTEM + 4)
CMD_QUERY_ARUND_USERS =(CMD_SOCIAL_SYSTEM + 5)
CMD_SET_INTERES= (CMD_SOCIAL_SYSTEM + 6)
CMD_SET_BLACK=(CMD_SOCIAL_SYSTEM + 7)
CMD_QUERY_MY_USER_LIST  = (CMD_SOCIAL_SYSTEM + 8)
CMD_QUERY_MY_FANS = (CMD_SOCIAL_SYSTEM + 9)
CMD_QUERY_ALBUM = (CMD_SOCIAL_SYSTEM + 10)
CMD_ADD_TO_ALBUM = (CMD_SOCIAL_SYSTEM + 11)
CMD_SEND_FLOWER = (CMD_SOCIAL_SYSTEM + 12)
CMD_REPORT = (CMD_SOCIAL_SYSTEM + 13)
CMD_QUERY_USER_INFO_BYSID =  (CMD_SOCIAL_SYSTEM + 15) #//通过sid 查询用户信息
CMD_ENFORCE_PHOTO = (CMD_SOCIAL_SYSTEM + 16)  #//强看私照

#共享私照模块
#define CMD_INVATE_FOR_SHARE (CMD_SOCIAL_SYSTEM + 101)
#define CMD_RESPONSE_FOR_SHARE (CMD_SOCIAL_SYSTEM + 102)
CMD_RELEASE_SHARE = (CMD_SOCIAL_SYSTEM + 103)

CMD_VALUE_PHOTO =  (CMD_SOCIAL_SYSTEM + 301) #照片评分

CMD_NEWS_LIST = (CMD_SOCIAL_SYSTEM + 401)  #获取广场信息列表
#define CMD_QUERY_NEWS_LIST  CMD_NEWS_LIST
CMD_UPDATE_NEWS =  (CMD_SOCIAL_SYSTEM + 402) #更新我的广场信息
CMD_CANCLE_NEWS = (CMD_SOCIAL_SYSTEM + 403) #下架我的广场信息
CMD_MY_NEWS_LIST = (CMD_SOCIAL_SYSTEM + 404)  #获取我的广场信息列表

#摇一摇
#define CMD_SHAKE_SHAKE (CMD_SOCIAL_SYSTEM + 104)
#偶遇模块
CMD_REPORT_LOCATION = (CMD_SOCIAL_SYSTEM + 201)

#IM系统
#请求/应答包
CMD_SEND_CHAT_MSG  = (CMD_IM_SYSTEM + 1)
CMD_PRIVATE_MSG_RECEIVED = (CMD_IM_SYSTEM + 2)
#define CMD_SYSTEM_MSG_RECEIVED (CMD_IM_SYSTEM + 3)
CMD_HELLO =  (CMD_IM_SYSTEM + 4)
#推送包
#define CMD_CHAT_MSG_NOTIFICATION (CMD_SOCIAL_SYSTEM + CMD_PUSH + 1)
#define CMD_SYSTEM_MSG_NOTIFICATION (CMD_SOCIAL_SYSTEM + CMD_PUSH + 2)

CMD_APPCFG_SYSTEM = CMD_APP_SYSTEM +1

ECMD_VERSION = 2

# SVR_IP = "120.25.239.220"
SVR_IP = "127.0.0.1"
SVR_PORT = 4000

# CONN_IP = "120.25.239.220"
CONN_IP = "127.0.0.1"
CONN_PORT = 9100
TEST_PHONE_NUM = '+8613812341111'
# TEST_PHONE_NUM = '+8618816822228'
TEST_PWD = '25d55ad283aa400af464c76d713c07ad'
TEST_UID = 100030
TEST_SID = 10007

MY_IP = "127.0.0.1"
MY_IP_N = socket.inet_aton(MY_IP)

MY_PORT = 29543
CON_LEN = 24
SVRADD = (SVR_IP, SVR_PORT)

DBHEADRES="00000000000000000000000000"
toRes="aaaaa"
fromRes="a    "
firstUin = 40201

class TestInfo:
    def __init__(self, name, passwd, level):
        self.name = name
        self.passwd = passwd
        self.prilevel= level

def OnlyTcpSendReq(cmd, PkgBody , sock, Uid=0, Sid = 0):
    ALL_len = CON_LEN + len(PkgBody)
    #ALL_len =0
    print "SID:", Sid, ", Uid =" , Uid
    PKG_HEAD = struct.pack("!HHHHIQI", ALL_len, cmd, 1, 1234, Sid, Uid, 0)
    #PKG_DB_HEAD = struct.pack("!IHIHI5sI5sII26s",0,0,0,0,Uin, fromRes, 0,toRes,0, 0,  DBHEADRES)
    PKG=PKG_HEAD+PkgBody
    PKGALL =PKG
    sock.send(PKGALL)
    print "req:", PkgBody
    
        
def TcpSendReq2(sock, cmd=0, req_json="", uid=0, sid=0):
    jsonStr = req_json
    fmt = '!%ds' % len(jsonStr)
    print 'input', jsonStr
    str = struct.pack(fmt, jsonStr)

    data = TcpSendReq(cmd, str, sock, uid, sid)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead = tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson) = struct.unpack(fmt, pkghead)
    
    return (errCode, strJson)

  
def OnlyTcpSendReq2(sock, cmd=0, req_json="", uid=0, sid=0):
    jsonStr = req_json

    fmt = '!%ds' % len(jsonStr)
    print 'OnlyTcpSendReq2() input', jsonStr

    str = struct.pack(fmt, jsonStr)

    #data  = TcpSendReq(cmd, str, sock, uid, sid)
    ALL_len = CON_LEN + len(str)
    #ALL_len =0
    print "OnlyTcpSendReq2 SID:", sid, ", Uid =", uid

    PKG_HEAD = struct.pack("!HHHHIQI", ALL_len, cmd, 1, 1234, sid, uid, 0)
    #PKG_DB_HEAD = struct.pack("!IHIHI5sI5sII26s",0,0,0,0,Uin, fromRes, 0,toRes,0, 0,  DBHEADRES)
    PKGALL=PKG_HEAD+str

    sock.send(PKGALL)


def TcpHello(uid, sid, sock):
    remarks = '{"platform":"%s","uid":"%d", "sid":"%d"}' % ("a", uid, sid)

    # fmt = '!%ds' % len(remarks)
    # str = struct.pack(fmt, remarks)
    #
    # data = TcpSendReq(10101, str, sock, 0)

    data = TcpSendReq(10100, remarks, sock, uid, sid)

    print 'TcpHello Received', data

    # tempstr = StringIO.StringIO(data)
    # tempstr.seek(0, 0)
    #
    # pkghead = tempstr.read(len(data))
    #
    # fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    # (Len, cmd, version, Seq , errCode, Uid, compress, strjson) = struct.unpack(fmt, pkghead)
    #
    # return (errCode, strjson)


def TcpSendReq(cmd, PkgBody, sock, Uid=0, Sid=0):
    ALL_len = CON_LEN + len(PkgBody)
    #ALL_len =0
    print "TcpSendReq Sid=%s, Uid=%s" % (Sid, Uid)
    PKG_HEAD = struct.pack("!HHHHIQI", ALL_len, cmd, 1, 1234, Sid, Uid, 0)
    #PKG_DB_HEAD = struct.pack("!IHIHI5sI5sII26s",0,0,0,0,Uin, fromRes, 0,toRes,0, 0,  DBHEADRES)
    PKG = PKG_HEAD + PkgBody
    PKGALL = PKG

    sock.send(PKGALL)

    print "send cmd=", cmd
    print "send PKGALL=", PKGALL

    sock.settimeout(2)
    try:
        print "recv"
        data = sock.recv(10240)
        print data
        return data
    except socket.error, msg:
        print msg
        sock.close()


def TcpRegister(sock, platform='a', phonenum='', passwd=''):
    jsonStr = '{"platform":"%s","phonenum":"%s","password":"%s"}' % (platform, phonenum, passwd)

    print "TcpRegister jsonStr=", jsonStr

    fmt = '!%ds' % len(jsonStr)
    str = struct.pack(fmt, jsonStr)

    data = TcpSendReq(10101, str, sock, 0)

    print 'Received', data
    
    if data == None:
        return 10000, "recv error"

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead = tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson) = struct.unpack(fmt, pkghead)

    return (errCode, strJson)


def default_login(s):
   return TcpLogin(s, platform='a', phonenum=TEST_PHONE_NUM, passwd=TEST_PWD)

g_uid = 0
def TcpLogin(sock, platform="a", phonenum="", cid="test", passwd=""):
    jsonStr = '{"platform":"%s","phonenum":"%s","password":"%s"}' % (platform, phonenum, passwd)
    fmt = '!%ds' % len(jsonStr)
    str= struct.pack(fmt, jsonStr)

    data = TcpSendReq(10102, str, sock, 0)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead = tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq, errCode, uid, compress, strJson) = struct.unpack(fmt, pkghead)
    return (errCode, strJson)

def TcpLogin2(sock, platform="a", uid=3000071, passwd="1xxswebd"):
    jsonStr = '{"uid":%d,"platform":"%s","password":"%s"}' % (uid, platform, passwd)
    fmt = '!%ds' % len(jsonStr)
    str= struct.pack(fmt, jsonStr)

    data = TcpSendReq(10102, str, sock, 0)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead = tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq, errCode, uid, compress, strJson) = struct.unpack(fmt, pkghead)
    return (errCode, strJson)


def TcpResetPwd(sock, uid=0, sid=0, oldpwd="", newpwd=""):
    jsonStr = '{"oldpsw":"%s", "newpsw":"%s"}' % (oldpwd, newpwd)
    fmt = '!%ds' % len(jsonStr)
    str = struct.pack(fmt, jsonStr)

    print "jsonStr", jsonStr

    data = TcpSendReq(10107, str, sock, uid, sid)

    print 'Received', data

    # tempstr = StringIO.StringIO(data)
    # tempstr.seek(0, 0)
    #
    # pkghead = tempstr.read(len(data))
    #
    # fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    # (Len, cmd, version, Seq, errCode, _, compress, strjson) = struct.unpack(fmt, pkghead)
    # return (errCode, strjson)

def TcpRetrievePwd(sock, uid=0, sid=0, phonenum="", password=""):
    jsonStr = '{"phonenum":"%s","password":"%s"}' % (phonenum, password)
    fmt = '!%ds' % len(jsonStr)
    str = struct.pack(fmt, jsonStr)

    data = TcpSendReq(10109, str, sock, uid, sid)

    print 'Received', data

def TcpBindPhone(sock, uid=0, sid=0, phonenum=""):
    jsonStr = '{"phonenum":"%s"}' % (phonenum)
    fmt = '!%ds' % len(jsonStr)
    str = struct.pack(fmt, jsonStr)

    data = TcpSendReq(10108, str, sock, uid, sid)

    print 'Received', data

def TcpSetDevToken(uid, token, sid, sock):
    jsonStr = '{"devicetoken":"%s"}' % (token)
    fmt = '!%ds' % len(jsonStr)
    str= struct.pack(fmt,jsonStr)

    data = TcpSendReq(10103, str, sock, uid, sid)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead = tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson)= struct.unpack(fmt, pkghead)
    return (errCode, strJson)


def TcpSetUserInfo(uid, sid, sock, did="", baseinfo="", exinfo=""):
    jsonStr = '{"did":"%s", "baseinfo":"%s", "exinfo":"%s"}' % (did,baseinfo,exinfo)
    fmt = '!%ds' % len(jsonStr)
    str = struct.pack(fmt,jsonStr)

    data = TcpSendReq(10104, str, sock,uid,sid)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead = tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson)= struct.unpack(fmt, pkghead)
    return (errCode, strJson)


def TcpGetUserInfo(uidlist, propertylist, sid, sock):
    jsonStr = '{"uidlist":['
    for k, uid in enumerate(uidlist):
        jsonStr=jsonStr + '%d,' % (uid)
    jsonStr=jsonStr.rstrip(',')
    jsonStr=jsonStr+'], "propertylist":['
    for k, v in enumerate(propertylist):
        jsonStr=jsonStr + '"%s",' % (v)
    jsonStr=jsonStr.rstrip(',')
    jsonStr=jsonStr+']}'
    fmt = '!%ds' % len(jsonStr)
    str= struct.pack(fmt,jsonStr)

    data  = TcpSendReq(10105, str, sock, uid, sid)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead =  tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson)= struct.unpack(fmt, pkghead)
    return (errCode, strJson)

def TcpGetUid(sock, phonenum="", uid=0, sid=0):
    jsonStr = '{"phonenum":"%s"}' % (phonenum)

    fmt = '!%ds' % len(jsonStr)
    str = struct.pack(fmt, jsonStr)

    data = TcpSendReq(10106, str, sock, uid, sid)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead =  tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson)= struct.unpack(fmt, pkghead)

    return (errCode, strJson)

# def TcpGetUid(uid, appkey, cidlist, sid, sock):
#     jsonStr = '{"appkey":"%s", "cidlist":[' % (appkey)
#     for k,v in enumerate(cidlist):
#         jsonStr=jsonStr + '"%s",' % (v)
#     jsonStr=jsonStr.rstrip(',')
#     jsonStr=jsonStr+']}'
#     fmt = '!%ds' % len(jsonStr)
#     str = struct.pack(fmt,jsonStr)
#
#     data = TcpSendReq(10106, str, sock, uid, sid)
#     print 'Received', data
#
#     tempstr = StringIO.StringIO(data)
#     tempstr.seek(0, 0)
#
#     pkghead = tempstr.read(len(data))
#
#     fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
#     (Len, cmd, version, Seq , errCode, uid, compress, strjson)= struct.unpack(fmt, pkghead)
#     return (errCode, strjson)

def TcpUserLocation(uid, xpos, ypos, sid, sock, level, hour=0, page=0):
    jsonStr = '{"Xpos":%f,"Ypos":%f,"Level":%d,"Hour":%d,"Page":%d}' % (xpos,ypos,level,hour,page)
    fmt = '!%ds' % len(jsonStr)
    print 'input', jsonStr
    str= struct.pack(fmt,jsonStr)

    data  = TcpSendReq(40001, str, sock, uid, sid)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead =  tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson)= struct.unpack(fmt, pkghead)
    return (errCode, strJson)

def TcpSendUserMsg(uid, sid, sock, content, touid, msgtype=0):
    jsonStr = '{"msgcontent":"%s","touid":%d,"msgtype":%d}' % (content,touid,msgtype)
    fmt = '!%ds' % len(jsonStr)
    print 'input', jsonStr
    str= struct.pack(fmt,jsonStr)

    data  = TcpSendReq(30101, str, sock, uid, sid)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead =  tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson)= struct.unpack(fmt, pkghead)
    return (errCode, strJson)

def TcpUserMsgReceived(uid, sid, sock, msgid, fromuid):
    jsonStr = '{"msgid":%d,"uid":%d}' % (msgid, fromuid)
    fmt = '!%ds' % len(jsonStr)
    print 'input', jsonStr
    str = struct.pack(fmt, jsonStr)

    data = TcpSendReq(30201, str, sock, uid, sid)
    print 'Received', data

    tempstr = StringIO.StringIO(data)
    tempstr.seek(0, 0)

    pkghead =  tempstr.read(len(data))

    fmt = "!HHHHIQI%ds" % (len(data)-CON_LEN)
    (Len, cmd, version, Seq , errCode, uid, compress, strJson)= struct.unpack(fmt, pkghead)
    return (errCode, strJson)


REGSVR_IP = "192.168.254.246"
REGSVR_PORT = 6450

def PushSend(cmd, PkgBody , Uin=0, Sid = 0):
    ALL_len = len(PkgBody)
    PKG = PkgBody
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    #s.connect((REGSVR_IP, REGSVR_PORT))
    s.sendto(PKG,(REGSVR_IP, REGSVR_PORT))
    s.settimeout(5)

    try:
        #s = socket.socket(af, socktype, proto)
        print s
    except socket.error, msg:
        s = None
        s.close()
        data,ADDR = s.recvfrom(4096)
        s.close()
        return data


