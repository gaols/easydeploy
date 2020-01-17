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

func TestDeployer(t *testing.T) {
	deployer := &Deployer{
		SrvConf: []*ServerConfig{
			{
				User:   "root",
				Server: "192.168.2.24",
			},
		},
	}

	deployer.Local("cd /home/gaols;ls")
	deployer.Upload("/home/gaols/Codes/go/src/github.com/gaols/easydeploy", "/tmp")
	deployer.Remote("cd /tmp/easydeploy;ls")
	deployer.OnceDoneDeploy(func(isDeployOk bool) error {
		return nil
	})
	deployer.Start()
}
