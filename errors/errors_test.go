package errors

import (
	"fmt"
	"testing"
)

func TestErrors_Trace(t *testing.T) {
	fmt.Println(ErrConfigHandlerNotImplement.Trace("name", "aa"))
}

func TestErrors_Trace2(t *testing.T) {
	err := New("test error", "name", "aa")
	fmt.Println(err.Trace("bb", "cc"))
}
