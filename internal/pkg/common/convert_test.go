package common_test

import (
	"harmoni/internal/pkg/common"
	"testing"
)

func TestBytesToString(t *testing.T) {
	s := "hello,world"
	b := common.StringToBytes(s)
	ss := common.BytesToString(b)
	t.Log(ss)
	if s != ss {
		t.FailNow()
	}
}
