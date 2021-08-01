package etcd

import (
	"testing"
)

func Test_intranetIP(t *testing.T) {

	gotIps, err := intranetIP()
	if err != nil {
		t.Errorf("intranetIP() error = %v", err)
		return
	}
	t.Log(gotIps)
}

