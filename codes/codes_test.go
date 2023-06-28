package codes

import (
	"fmt"
	"testing"
)

func TestCodes_Trace(t *testing.T) {
	fmt.Println(ErrInvalidArgs.Trace("name", "ljw"))
}

func TestErrors_Is1(t *testing.T) {
	arr := []struct {
		err     error
		checker func(err error) bool
		expect  bool
	}{
		{
			ErrInvalidPassword.Trace("name", "aa"),
			func(err error) bool {
				return Is(err, ErrInvalidPassword)
			},
			true,
		},
		{
			Trace(ErrInvalidPassword.Trace("name", "aa")),
			func(err error) bool {
				return Is(err, ErrInvalidPassword)
			},
			true,
		},
		{
			Trace(ErrInvalidPassword.Trace("name", "aa")),
			func(err error) bool {
				return Is(err, ErrPayOrderStatusInvalid)
			},
			false,
		},
		{
			Trace(ErrInvalidPassword.Trace("name", "aa")),
			func(err error) bool {
				return Is(err, New(CodeInvalidPassword, ""))
			},
			true,
		},
	}
	for index, v := range arr {
		if v.checker(v.err) != v.expect {
			t.Fatalf("check err(%d) fail, expect:%v got:%v", index, v.expect, v.checker(v.err))
		}
	}
}
