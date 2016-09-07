package common

//消息类型
const (
	MSG_SYS = iota
	MSG_USER
	MSG_TEAM

	MSG_NONE
)

//Msg Send
type UserMsgItem struct {
	CmdType   uint64 `json:"cmdtype"`
	MsgId     uint64 `json:"msgid"`
	FromUid   uint64 `json:"fromuid"`
	ToUid     uint64 `json:"touid"`
	Type      uint16 `json:"msgtype"`
	Content   string `json:"msgcontent"`
	SendTime  uint32 `json:"sendtime"`
	ApnsText  string `json:"apnstext,omitempty"`
	FType     int    `json:"ftype,omitempty"`
	FBv       int    `json:"frombv,omitempty"`
	ExtraData string `json:"extraData"`
}

type CSMsgItem struct {
	UserMsgItem
	FType int `json:"ftype,omitempty"`
}

//Msg Received
type MsgRecvedUser struct {
	MsgId uint64 `json:"msgid,omitempty"`
	Uid   uint64 `json:"uid,omitempty"`
}

//For Portal
type MsgSendItem struct {
	MsgId    uint64 `json:"msg_id"`
	FromUid  uint64 `json:"from_uid"`
	TouId    uint64 `json:"to_uid"`
	Content  string `json:"content"`
	SendTime uint32 `json:"send_time"`
	Type     int    `json:"type"`
}

type MsgRecvItem struct {
	MsgId    uint64 `json:"msg_id"`
	Uid      uint64 `json:"uid"`
	RecvTime uint32 `json:"recv_time"`
	Type     int    `json:"type"`
}
