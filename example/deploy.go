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
	deployer.Upload("/home/gaols/Codes/go/src/github.com/gaols/easydeploy", "/tmp/") // you can upload any file to remote server
	deployer.Remote("ps aufx") // you can run any shell command on remote server
	deployer.OnceDoneDeploy(func(deployOk bool) error { // you can do some clean after deployed
		_, err := easyssh.Local("ls -l /tmp")
		return err
	})
	deployer.Start() // start deploy
}
