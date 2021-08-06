package message

import "testing"

func TestSendMsg(t *testing.T) {
	if err := SendMsg("18826073368", "133333"); err != nil {
		t.Errorf("SendMsg() error = %v", err)
	}
}
