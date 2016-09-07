package syslog

import (
	"testing"
)

func TestSysLog(t *testing.T) {
	if err := SysLogInit("127.0.0.1:4150", "sysLogTopic"); err != nil {
		t.Fatal(err)
	}

	SysLog("this is a test syslog!!!")

	SysLogDeinit()
}
