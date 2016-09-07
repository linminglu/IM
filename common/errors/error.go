package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

func test() error {
	return errors.New("test")
}

type ErrorT struct {
	Code   string   `json:"code"`
	Caller []string `json:"caller"`
	Reason []string `json:"reason"`
}

func caller(dept int) string {
	_, file, line, ok := runtime.Caller(dept)
	if !ok {
		return "Unknown"
	}

	idx := strings.LastIndex(file, "/")

	return fmt.Sprint(file[idx+1:], ":", line)

}

// 新建错误
func New(code string, reason ...interface{}) *ErrorT {
	return &ErrorT{
		Code:   code,
		Caller: []string{caller(2)},
		Reason: []string{fmt.Sprint(reason)},
	}
}

// 追加错误
func (e *ErrorT) As(reason ...interface{}) *ErrorT {
	e.Caller = append(e.Caller, caller(2))
	e.Reason = append(e.Reason, fmt.Sprint(reason))
	return e
}

// 复制
func (e *ErrorT) Clone() *ErrorT {
	return &ErrorT{
		Code:   e.Code,
		Caller: append([]string{}, e.Caller...),
		Reason: append([]string{}, e.Reason...),
	}
}

// error 接口实现
func (e *ErrorT) Error() string {
	count := len(e.Caller)
	errStr := fmt.Sprintln("err_code:", e.Code)
	for i := 0; i < count; i++ {
		errStr += fmt.Sprintf("err_stack_%v==> caller:%v reason:%v\n", i, e.Caller[i], e.Reason[i])
	}
	return errStr
}

// err 是否相同
func (e *ErrorT) Equal(err *ErrorT) bool {
	if e.Code == err.Code {
		return true
	}
	return false
}
