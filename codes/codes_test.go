package codes

import (
	"fmt"
	"testing"
)

func TestCodes_Trace(t *testing.T) {
	fmt.Println(ErrInvalidArgs.Trace("name", "ljw"))
}
