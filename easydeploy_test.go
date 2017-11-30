package easydeploy

import (
	"testing"
)

func TestNewSrvConf(t *testing.T) {
	sc := NewSrvConf("gaols@192.168.1.100:22/password@**/")
	if sc.Password != "password@**/" {
		t.Error("parse password failed")
	}
	if sc.User != "gaols" {
		t.Error("parse user failed")
	}
	if sc.Server != "192.168.1.100" {
		t.Error("parse server failed")
	}
	if sc.Port != "22" {
		t.Error("parse port failed")
	}
	t.Log(sc)
}
