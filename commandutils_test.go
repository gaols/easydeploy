package easydeploy

import (
	"testing"
	"github.com/gaols/easyssh"
)

func TestTar(t *testing.T) {
	cmd := Tar("/tmp/easydeploy.tar.gz", "/home/gaols/Codes/go/src/github.com/gaols/easydeploy/")
	cmd.Run(nil, nil)
	if !easyssh.IsRegular("/tmp/easydeploy.tar.gz") {
		t.Error("tar command failed")
	}
	defer easyssh.Local("cd /tmp;rm -f easydeploy.tar.gz")
}

func TestUnTar(t *testing.T) {
	cmd := Tar("/tmp/easydeploy.tar.gz", "/home/gaols/Codes/go/src/github.com/gaols/easydeploy/")
	cmd.Run(nil, nil)
	cmd = UnTar("/tmp/easydeploy.tar.gz", "/tmp/")
	cmd.Run(nil, nil)
	defer easyssh.Local("cd /tmp;rm -f easydeploy.tar.gz")
	defer easyssh.Local("cd /tmp;rm -rf easydeploy")

	if !easyssh.IsDir("/tmp/easydeploy") {
		t.Error("untar command failed")
	}

}
