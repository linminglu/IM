package errors

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("100", "test", "test1")
	fmt.Println(err)

	err.As("test2")
	fmt.Println(err.Error())
}
