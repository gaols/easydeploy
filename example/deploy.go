package main

import "github.com/gaols/easydeploy"

func main() {
	deployer := &easydeploy.Deployer{
		SrvConf: []*easydeploy.ServerConfig{
			{
				User:     "gaols",
				Server:   "192.168.2.155",
				Port:     "22",
				Password: "gaolsz",
			},
		},
	}
	deployer.Upload("/home/gaols/Codes/go/src/github.com/gaols/easydeploy", "/tmp/")
	deployer.Remote("ps aufx")
	deployer.OnceDoneDeploy(func(deployOk bool) error {
		return nil
	})
	deployer.Verbose()
	deployer.Start()
}
