# easydeploy

## Deployment

To my opinion, a simple deployment just the following four steps:

* prepare the artifacts to be deployed
* put the prepared artifacts to your deploy servers
* restart deploy servers
* do some clean

Now with easydeploy you can:

```go
import (
	"github.com/gaols/easydeploy"
	"github.com/gaols/easyssh"
)

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
    _, err := easyssh.Local("ls -l /home/tmp")
    return err
})
deployer.Verbose()
deployer.Start()
```

## Notes

1. `Upload(localPath, remotePath string)` upload a local file or local dir to its corresponding remote path, the remotePath 
should contain the file name if the localPath is a regular file, however, if the localPath to copy is dir, the remotePath must
be the dir into which the localPath will be copied.
2. The local commands and remote commands you registered by calling `Local/Remote/Upload` to deployer will not run until
`Start()` method being called. 
3. `Local/Remote/Upload` method can be called multiple times to fit your deployment needs.  
4. If you'd like to run the local or remote shell command manually, please refer to [easyssh](https://github.com/gaols/easyssh).

## So easy to deploy

[A simple deploy sample](https://github.com/gaols/easydeploy/blob/master/example/deploy.go)
