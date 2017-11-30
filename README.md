# easydeploy

## Description

Package easydeploy makes your deployment easy.

## Deployment

To my opinion, a simple deployment just the following four steps:

* prepare the artifacts to be deployed
* put the prepared artifacts to your deploy servers
* restart servers
* do some clean

Now with easydeploy you can:

```
deployer := &easydeploy.Deployer{
    SrvConf: []*easydeploy.ServerConfig{
        easydeploy.NewSrvConf("gaols@192.168.1.100:22/123456"),
        easydeploy.NewSrvConf("gaols@192.168.1.103:22/123456"),
    },
}

// step 1
deployer.Local("/path/to/your/prepare-artifacts.sh")
// step 2
deployer.Upload("/path/to/your/artifacts", "/path/to/remote")
// step 3
deployer.Remote("/path/to/your/restart-server-on-remote.sh")
// step 4
deployer.OnceDoneDeploy(func(isDeployOk bool) error {
    return nil
})
deployer.Verbose()
deployer.Start()
```
