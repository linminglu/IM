package main

const (
	MAX_RECVBUF_SIZE = 2048

	NET_READ_DEADLINE  = 10
	NET_WRITE_DEADLINE = 10

	NET_DISCONNECT_DEADLINE = 120 // 两分钟无操作断开
)

const (
	IM_ERR_SUCCESS           = 0
	IM_ERR_NET_TIMEOUT       = 1
	IM_ERR_NETWORK           = 2
	IM_ERR_DATA_PACKAGE_SIZE = 3
	IM_ERR_DATA_FORMAT       = 4
	IM_ERR_SYSTEM            = 5
)

// protpcol
const (
	CommandKeyword = "command"
	//command
	IM_CMD_BEATHEART = 1000
)
