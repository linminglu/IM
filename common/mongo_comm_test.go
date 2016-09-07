package common

import (
	"fmt"
	"testing"
)

func TestGetUserMsg(t *testing.T) {
	ret := MongoInit("192.168.20.51:27017")
	if ret != 0 {
		fmt.Println("MongoInit error")
		return
	}

	userMsg := &UserMsgItem{}
	msgs, err := userMsg.GetUserMsg(100057)

	if err != nil {
		fmt.Println("GetUserMsg error:", err)
		return
	}

	for _, msg := range msgs {
		fmt.Println(*msg)
	}

}
