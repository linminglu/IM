package syslog_server

import (
	"testing"
)

func TestFile(t *testing.T) {
	log := NewLogFile(".", "log_test_file")
	log.WriteString("test")
	log.date = log.date.AddDate(0, 0, -1)
	log.WriteString("test")
}
