package tcp

import (
	"testing"
)

// go test -v -test.run Test_GetMsgId
func Test_GetMsgId(t *testing.T) {
	if 2 != getMsgId() {
		t.Errorf("getMsgId fail")
	}
	if 3 != getMsgId() {
		t.Errorf("getMsgId fail")
	}
	// for test
	globalMsgId = MaxInt64 - 3
	for {
		msgId := getMsgId()
		if msgId == MaxInt64 {
			if 2 != getMsgId() {
				t.Errorf("getMsgId fail")
				return
			}
			break
		}
	}
}
