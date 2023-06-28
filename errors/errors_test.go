package errors

import (
	"fmt"
	"testing"
)

func TestErrors_Trace(t *testing.T) {
	fmt.Println(ErrConfigHandlerNotImplement.Trace("name", "aa"))
}

func TestErrors_Trace2(t *testing.T) {
	err := New("test error", "k1", "v1")
	fmt.Println(err.Trace("k2", "v2"))
}

func TestErrors_Trace3(t *testing.T) {
	err := New("test error", "k1", "v1").Trace("k2", "v2")
	fmt.Println(Trace(err, "k3", "v3"))
}

func TestErrors_Is1(t *testing.T) {
	arr := []struct {
		err     error
		checker func(err error) bool
		expect  bool
	}{
		{
			ErrConfigHandlerNotImplement.Trace("name", "aa"),
			func(err error) bool {
				return Is(err, ErrConfigHandlerNotImplement)
			},
			true,
		},
		{
			Trace(ErrConfigHandlerNotImplement.Trace("name", "aa")),
			func(err error) bool {
				return Is(err, ErrConfigHandlerNotImplement)
			},
			true,
		},
		{
			Trace(ErrConfigHandlerNotImplement.Trace("name", "aa")),
			func(err error) bool {
				return Is(err, ErrAdminClientNotImplement)
			},
			false,
		},
	}
	for index, v := range arr {
		if v.checker(v.err) != v.expect {
			t.Fatalf("check err(%d) fail, expect:%v got:%v", index, v.expect, v.checker(v.err))
		}
	}
}
