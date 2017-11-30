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
			NewSrvConf("gaols@192.168.1.100:22/123456"),
		},
	}

	deployer.Local("/path/to/your/prepare-artifacts.sh")
	deployer.Upload("/path/to/your/artifacts", "/path/to/remote")
	deployer.Remote("/path/to/your/restart-server-on-remote.sh")
	deployer.OnceDoneDeploy(func(isDeployOk bool) error {
		return nil
	})
	deployer.Start()
}
