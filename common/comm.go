package common

import (
	//	//#include "../cpp/status_cache.c"
	//	"C"
	//	"unsafe"

	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/donnie4w/go-logger/logger"
)

//消息头部
type PkgHead struct {
	PkgLen uint16 //整个包的长度
	Cmd    uint16 //命令字
	Ver    uint16 //协议版本号
	Seq    uint16 //通讯流水 ID
	Sid    uint32 //Session ID
	Uid    uint64 //UID
	Flag   uint32 //包体压缩 & 包体加密
}

func (head PkgHead) ToString() string {
	str, err := json.Marshal(head)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return string(str)
}

//消息尾部
type InnerPkgTail struct {
	ConnIP   int64  //连接服务器IP
	ConnPort uint32 //连接服务器PORT
	FromUid  uint64 //消息来源用户ID
	ToUid    uint64 //消息去向用户ID
	Sid      uint32 //Session ID
	MsgId    uint64 //Msg ID
}

func (tail InnerPkgTail) ToString() string {
	str, err := json.Marshal(tail)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return string(str)
}

type TextTeamMsg struct {
	FromUid    uint64 `json:"fromuid"`
	ToTeamId   uint64 `json:"toteamid"`
	SendTime   int    `json:"sendtime"`
	MsgContent string `json:"msgcontent"`
	MsgId      uint64 `json:"msgid"`
	MsgType    int    `json:"msgtype"`
	ApnsText   string `json:"apnstext,omitempty"`
	FBv        int    `json:"frombv"`
}

type TextChatMsg struct {
	FromUid    uint64 `json:"fromuid"`
	ToUid      uint64 `json:"touid"`
	SendTime   int    `json:"sendtime"`
	MsgContent string `json:"msgcontent"`
	MsgId      uint64 `json:"msgid"`
	MsgType    int    `json:"msgtype"`
}

type TxtMsg struct {
	TextChatMsg TextChatMsg `json:"text_chat_msg,omitempty"`
}

type SystemMsg struct {
	MsgType      int    `json:"msgtype"`
	SysTeamMsgId uint64 `json:"systemmsgid"`
	MsgContent   string `json:"msgcontent"`
	SendTime     int    `json:"sendtime"`
	ToUid        uint64 `json:"touid"`
}

type UserMsgSave struct {
	MsgId    uint64 `json:"msgid"`
	FromUid  uint64 `json:"fromuid"`
	ToUid    uint64 `json:"touid"`
	Type     uint16 `json:"msgtype"`
	MsgBuff  []byte `json:"msgcontent"`
	SendTime uint32 `json:"sendtime"`
}

type TeamMsgSaveItem struct {
	Opt      int32 //0-delete 1-insert
	MsgId    uint64
	TouId    uint64
	TeamId   uint64
	Msg      []byte
	SendTime uint32
}

type UserInfo struct {
	Uid      uint64 `json:"uid"`
	Pwd      string `json:"passwd"`
	PhoneNum string `json:"phonenum"`
	Platform string `json:"platform"`
	Did      string `json:"did"`
	BaseInfo string `json:"baseinfo"`
	ExInfo   string `json:"exinfo"`
	RegDate  uint64 `json:"regdate"`
	BV       uint64 `json:"bv"`
	V        uint64 `json:"v"`
}

type TeamInfo struct {
	Uid      uint64 `json:"uid,omitempty"`
	TeamId   uint64 `json:"teamid,omitempty"`
	TeamType int    `json:"teamtype,omitempty"`
	TeamName string `json:"teamname,omitempty"`
	CoreInfo string `json:"coreinfo,omitempty"`
	ExInfo   string `json:"exinfo,omitempty"`
	MaxCount int    `json:"maxcount,omitempty"`
	IV       int64  `json:"infov,omitempty"`
	MV       int64  `json:"memberv,omitempty"`
}

type ApnsMsg struct {
	//	AppKey string
	Msg     string
	UidList []uint64
}

type CSInfo struct {
	Uid     uint64
	Pwd     string
	AppKey  string
	Account string
	Name    string
	Image   string
	Tel     string
	Email   string
	Enable  int
	V       int
}

func (this *CSInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"uid":      this.Uid,
		"appkey":   this.AppKey,
		"cid":      this.Account,
		"nickname": this.Name,
		"fid":      this.Image,
		"tel":      this.Tel,
		"email":    this.Email,
		"enable":   this.Enable,
		"v":        this.V,
	})
}

type CSVerInfo struct {
	CsId uint64 `json:"csid,omitempty"`
	V    int    `json:"v,omitempty"`
}

const (
	ERROR_NONE         = iota // 无错误
	ERROR_UNKNOWN             // 未知错误
	ERROR_SERVER_BUSY         // 服务器繁忙
	ERROR_CLIENT_BUG          // 客户端请求数据包异常
	ERROR_OUT_OF_REACH        // 未能到达，权限不够（如金币不足，或无法查看私照）
	ERROR_TOUCH_TOP           // 已到达最大值（如周边用户已全部拉取完成）
	ERROR_NO_DATA             // 无可用数据返回
	ERROR_TIMEOUT             // 已过期
	ERROR_OUT                 // 要求离开
	ERROR_ACCOUNT             // 注册帐号已存在
	IN_BLACKLIST
	NOT_WHITELIST
	ERR_CID
	ERR_APPKEY             //13 AppKey错误
	ERR_AUTOLOGIN_CONFLICT //14 自动登录失败
	ERR_TOKENLOGIN_EXPIRED //15 token已过期
	ERR_TOKENLOGIN_UNKNOWN //17 token登录其它错误

	ERR_CODE_ERR_PKG = 3

	ERR_CODE_PASSWD        = 100
	ERR_CODE_NET           = 101
	ERR_CODE_SYS           = 102
	ERR_CODE_TIMEOUT       = 103
	ERR_CODE_NO_USER       = 104
	ERR_CODE_TRAM_MAXNUM   = 105
	ERR_CODE_TEAM_PRI      = 106
	ERR_CODE_CID_EXIST     = 108
	ERR_CODE_CS_ENABLE     = 109
	ERR_CODE_CS_LOGIN      = 110
	ERR_CODE_PHONENUM_USED = 111

	CONN_TIMEOUT         = 180
	CONN_NOLOGIN_TIMEOUT = 30

	MAX_TEAM_NUM_PER    = 300
	MAX_MEMBER_NUM_TRAM = 500
	MAX_TEAM_MSG_PER    = 50
	SIZEOF_PKGHEAD      = 24
	SIZEOF_INNERTAIL    = 40

	ERROR_RPC_SERVER_BUSY = `{"err_code":1001, "total": 0, "msg_id":0}`
	ERROR_RPC_VERCODE     = `{"err_code":1010, "total": 0, "msg_id":0}`
	ERROR_RPC_NOTARGET    = `{"err_code":1011, "total": 0, "msg_id":0}`
)

const (
	//用户系统 UserSystem
	DU_CMD_USER_HELLO               = 10100 // 心跳
	DU_CMD_USER_REGISTER            = 10101 // 注册或登录
	DU_CMD_USER_LOGIN               = 10102 // 登录
	DU_CMD_USER_SET_DEVICE_TOKEN    = 10103 // 设置DeviceToken，iOS专用
	DU_CMD_USER_SET_MY_INFO         = 10104 // 设置本用户信息
	DU_CMD_USER_GET_INFO            = 10105 // 请求或批量请求用户信息
	DU_CMD_USER_GET_UID             = 10106 // 查询用户标识
	DU_CMD_USER_RESET_PWD           = 10107 // reset pwd
	DU_CMD_USER_BIND_PHONE          = 10108 // 绑定手机号
	DU_CMD_USER_RETRIEVE_PWD        = 10109 // 找回密码
	DU_CMD_USER_PURE_REGISTERLOGIN  = 10111
	DU_CMD_USER_PURE_REGISTER       = 10111
	DU_CMD_USER_TOKEN_LOGIN         = 10112 // token登录
	DU_PUSH_CMD_USER_LOGIN_CONFLICT = 15100 // 别处登录通知
	DU_CMD_USER_GET_APP_INFO        = 11000
	DU_CMD_USER_GET_SETUP_ID        = 11001
	DU_CMD_USER_SET_SETUP_ID        = 11002
	DU_CMD_GET_CSID_LIST            = 11003
	DU_CMD_GET_CSINFO_LIST          = 11004

	//即时通信系统 IMSystem
	DU_CMD_IM_SEND_USER_MSG       = 30101 // 发送用户间IM消息
	DU_CMD_IM_SEND_TEAM_MSG       = 30102 // 发送小组IM消息
	DU_CMD_IM_SYSTEM_MSG_RECEIVED = 30200 // 系统消息已送达
	DU_CMD_IM_USER_MSG_RECEIVED   = 30201 // 用户间消息已送达
	DU_CMD_IM_TEAM_MSG_RECEIVED   = 30202 // 小组消息已送达
	DU_PUSH_CMD_IM_SYSTEM_MSG     = 35100 // 系统IM消息通知
	DU_PUSH_CMD_IM_USER_MSG       = 35101 // 用户IM消息通知
	DU_PUSH_CMD_IM_TEAM_MSG       = 35102 // 小组IM消息通知
	DU_PUSH_CMD_IM_REPORT_MSG     = 30100

	//周边系统 AroundSystem
	DU_CMD_AROUND_QUERY = 40001 //拉取周边用户位置信息

	//小组系统 TeamSystem
	DU_CMD_TEAM_CREATE        = 50001 // 创建新小组
	DU_CMD_TEAM_DELETE        = 50002 // 删除小组
	DU_CMD_TEAM_GET_INFO      = 50011 // 批量获取小组信息
	DU_CMD_TEAM_GET_ALL       = 50012 // 获取某个用户所有小组
	DU_CMD_TEAM_GET_SYS       = 50013 // 获取系统预设所有小组
	DU_CMD_TEAM_SET_INFO      = 50021 // 设置小组信息
	DU_CMD_TEAM_ADD_MEMBER    = 50022 // 小组里添加用户
	DU_CMD_TEAM_REMOVE_MEMBER = 50023 // 小组里删除用户
	DU_CMD_TEAM_GET_MEMBER    = 50024 // 查找
	DU_CMD_ADD_MEMBER2WB      = 50025 //黑白名单添加用户
	DU_CMD_DEL_MEMBER2WB      = 50026 //黑白名单删除用户
	DU_CMD_GET_MEMBER2WB      = 50027 //黑白名单查询用户
	DU_CMD_TEAM_END           = 50028

	EXCHANGE_REG = "ex-reg"

	PT_IOS     = 1
	PT_ANDROID = 2
	PT_WP      = 3
	PT_WEB     = 4
	PT_PC      = 5
	PT_KF      = 6

	MIN_SETUP_ID = 100000000000
)

func IntToIP(ip_int int64) string {
	result := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		result[3-i] = byte((ip_int >> uint(8*i)) & 0xff)
	}
	return result.String()
}

func InetAton(ip net.IP) uint64 {
	var sum uint64 = 0

	sum += uint64(ip[12]) << 24
	sum += uint64(ip[13]) << 16
	sum += uint64(ip[14]) << 8
	sum += uint64(ip[15])

	return sum
}

func daemon(nochdir, noclose int) int {
	var ret, ret2 uintptr
	var err syscall.Errno

	darwin := runtime.GOOS == "darwin"

	// already a daemon
	if syscall.Getppid() == 1 {
		return 0
	}

	// fork off the parent process
	ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		return -1
	}

	// failure
	if ret2 < 0 {
		os.Exit(-1)
	}

	// handle exception for darwin
	if darwin && ret2 == 1 {
		ret = 0
	}

	// if we got a good PID, then we call exit the parent process.
	if ret > 0 {
		os.Exit(0)
	}

	/* Change the file mode mask */
	_ = syscall.Umask(0)

	// create a new SID for the child process
	s_ret, s_errno := syscall.Setsid()
	if s_errno != nil {
		logger.Error("Error: syscall.Setsid errno: %d", s_errno)
	}
	if s_ret < 0 {
		return -1
	}

	if nochdir == 0 {
		os.Chdir("/")
	}

	if noclose == 0 {
		f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
		if e == nil {
			fd := f.Fd()
			syscall.Dup2(int(fd), int(os.Stdin.Fd()))
			syscall.Dup2(int(fd), int(os.Stdout.Fd()))
			syscall.Dup2(int(fd), int(os.Stderr.Fd()))
		}
	}

	return 0
}

//func StatCacheInit(shmkey uint32) {
//	C.UserStatCacheInit(C.uint(shmkey))
//}
//
////int UserStatCache_get(UINT64   uid, UINT64 * psid, UINT32  *pip, UINT16 * pport, UINT32 * plasttime )
//func StatCacheGet(pinfo *InnerPkgTail) int {
//	logger.Debug("----------StatCacheGet:pinfo.FromUid", pinfo.FromUid)
//
//	var uid uint64 = pinfo.FromUid
//	var sid uint64 = 0
//	var ip uint32 = 0
//	var port uint16 = 0
//	var lasttime uint32 = 0
//	var stat uint8 = 0
//	var rs int8 = 0
//	ret := C.UserStatCacheGet(
//		C.ulong(uid),
//		(*C.ulong)(unsafe.Pointer(&sid)),
//		(*C.uint)(unsafe.Pointer(&ip)),
//		(*C.ushort)(unsafe.Pointer(&port)),
//		(*C.uint)(unsafe.Pointer(&lasttime)),
//		(*C.uchar)(unsafe.Pointer(&stat)),
//		(*C.char)(unsafe.Pointer(&rs)))
//
//	logger.Debug("----------StatCacheGet:", "uid:", uid, "sid:", sid, "ip:", ip, "port:", port, "stat:", stat, "rs:", rs)
//
//	if ret == 0 {
//		pinfo.ConnIP = int64(ip)
//		pinfo.ConnPort = uint32(port)
//		pinfo.Sid = uint32(sid)
//		pinfo.ToUid = uint64(stat)
//		pinfo.MsgId = uint64(rs)
//		return 0
//	}
//
//	return 1
//}
//
//func StatCacheSet(info InnerPkgTail) {
//	logger.Debug("----------StatCacheSet info:", info)
//
//	var uid uint64 = info.FromUid
//	var sid uint64 = 0
//	var ip uint32 = 0
//	var port uint16 = 0
//	var lastTime uint32 = 0
//	var stat uint8 = 0
//	var rs int8 = 0
//
//	ret := C.UserStatCacheGet(C.ulong(uid),
//		(*C.ulong)(unsafe.Pointer(&sid)),
//		(*C.uint)(unsafe.Pointer(&ip)),
//		(*C.ushort)(unsafe.Pointer(&port)),
//		(*C.uint)(unsafe.Pointer(&lastTime)),
//		(*C.uchar)(unsafe.Pointer(&stat)),
//		(*C.char)(unsafe.Pointer(&rs)))
//
//	if info.ToUid == 0 && uint64(info.Sid) != sid {
//		logger.Debug("logout sync local sid ", sid, " != sync sid", info.Sid)
//		return
//	}
//
//	logger.Debug("set:", info.FromUid, info.Sid, info.ConnIP, info.ConnPort, info.ToUid, info.MsgId)
//
//	uid = uint64(info.FromUid)
//	sid = uint64(info.Sid)
//	ip = uint32(info.ConnIP)
//	port = uint16(info.ConnPort)
//	stat = uint8(info.ToUid)
//	rs = int8(info.MsgId)
//
//	ret = C.UserStatCacheSet(C.ulong(uid),
//		(C.ulong)(sid),
//		(C.uint)(ip),
//		(C.ushort)(port),
//		(C.uchar)(stat),
//		(C.char)(rs))
//
//	logger.Debug("StatCacheSet ", uid, sid, " ret = ", ret)
//}
//
////int UserStatCacheScan(const char * filename, long mintime, long maxtime)
//func StatCacheScan(filename string, mintime int64, maxtime int64) int {
//	if mintime > maxtime {
//		logger.Error("input error. ")
//		return 1
//	}
//	logger.Info("mintime: ", mintime, " maxtime: ", maxtime)
//
//	ret := C.UserStatCacheScan(C.CString(filename),C.long(mintime),C.long(maxtime))
//
//	logger.Debug("filename: ", filename, " ret = ", ret)
//
//	if ret == 0 {
//		logger.Info("Scan success.")
//		return 0
//	}
//
//	return 1
//}

func DecPkgBody(body []byte) (int, PkgHead) {
	p := bytes.NewReader(body)
	head := PkgHead{}

	err := binary.Read(p, binary.BigEndian, &head)
	if err != nil {
		logger.Error("read pkghead fail:", err.Error())
		return 1, head
	}

	return 0, head
}

func DecPkgInnerBody(body []byte) (int, PkgHead, []byte, InnerPkgTail) {
	p := bytes.NewReader(body)
	head := PkgHead{}
	tail := InnerPkgTail{}

	err := binary.Read(p, binary.BigEndian, &head)
	if err != nil {
		logger.Error("read pkghead fail:", err.Error())
		return 1, head, nil, tail
	}

	var jsonStr []byte

	if head.PkgLen > SIZEOF_PKGHEAD {
		jsonStr = make([]byte, head.PkgLen-SIZEOF_PKGHEAD)
		err = binary.Read(p, binary.BigEndian, &jsonStr)
		if err != nil {
			logger.Error("DecPkgInnerBody() read req json fail:", err.Error())
			return 2, head, nil, tail
		}
	} else {
		jsonStr = []byte("")
	}

	//	logger.Debug("DecPkgInnerBody() head.Uid:", head.Uid)
	//	logger.Debug("DecPkgInnerBody() resp json:", string(jsonStr[:]))

	err = binary.Read(p, binary.BigEndian, &tail)
	if err != nil {
		logger.Error("read req innertail fail:", err.Error())
		return 3, head, nil, tail
	}

	return 0, head, jsonStr, tail
}

//func MsgDBFlagSet(uid uint64) int {
//	C.MsgDBFlagSet(C.ulong(uid))
//	return 0
//}
//
//func MsgDBFlagClear(uid uint64) int {
//	C.MsgDBFlagClear(C.ulong(uid))
//	return 0
//}

//func MsgDBFlagGet(uid uint64) int {
//	ret := C.MsgDBFlagGet(C.ulong(uid))
//	return int(ret)
//}

const UID_FLAG = 0xffffffffff000000

func GetLongUid(shortuid, platform uint64) uint64 {
	uid := (platform & 0xf) + (shortuid << 24)
	return uid
}

//func GetKefuUid(appid uint64) uint64 {
//	uid := ((appid << 4) & 0xffffff)
//	return uid
//}

func IsKefuUid(uid uint64) bool {
	return ((uid & 0xf) == uint64(PT_KF))
}

func GetShortUid(longuid uint64) uint64 {
	return longuid >> 24
}

func GetAppIdPtFromUid(uid uint64) (uint64, uint64) {
	return ((uid & 0xffffff) >> 4), (uid & 0xf)
}

func GetChatMsgId() uint64 {
	return uint64(time.Now().UnixNano()) / 1000
}

func GetTeamMsgId() uint64 {
	//return uint64(time.Now().UnixNano() / 10) | 0x2000000000000000
	return uint64(time.Now().UnixNano()/1000) | 0x10000000000000
}

func GetSysMsgId() uint64 {
	//return uint64(time.Now().UnixNano() / 10) | 0x4000000000000000
	return uint64(time.Now().UnixNano()/1000) | 0x20000000000000
}

func IsChatMsgId(msgid uint64) bool {
	return (msgid < 0x10000000000000)
}

func IsTeamMsgId(msgid uint64) bool {
	//return uint64(time.Now().UnixNano() / 10) | 0x2000000000000000
	return (msgid & 0x10000000000000) > 0
}

func IsSysMsgId(msgid uint64) bool {
	//return uint64(time.Now().UnixNano() / 10) | 0x4000000000000000
	return (msgid & 0x20000000000000) > 0
}

func ConnectToNSQAndLookupd(r *nsq.Consumer, nsqAddrs []string, lookupd []string) error {
	for _, addrString := range nsqAddrs {
		logger.Info("add nsqd addr %s", addrString)
		err := r.ConnectToNSQD(addrString)
		if err != nil {
			return err
		}
	}

	for _, addrString := range lookupd {
		logger.Info("add lookupd addr %s", addrString)
		err := r.ConnectToNSQLookupd(addrString)
		if err != nil {
			return err
		}
	}

	return nil
}
