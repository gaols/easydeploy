package main

import (
	"github.com/gaols/easydeploy"
	"github.com/gaols/easyssh"
)

func main() {
	deployer := &easydeploy.Deployer{
		SrvConf: []*easydeploy.ServerConfig{
			{
				User:     "gaols",
				Server:   "192.168.2.155",
				Port:     "22",
				Password: "******",
			},
		},
	}
	deployer.Upload("/home/gaols/Codes/go/src/github.com/gaols/easydeploy", "/tmp/")
	deployer.Remote("ps aufx")
	deployer.OnceDoneDeploy(func(deployOk bool) error {
		_, err := easyssh.Local("ls -l /home/tmp")
		return err
	})
	deployer.Verbose()
	deployer.Start()
}
