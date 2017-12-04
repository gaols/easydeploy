package easydeploy

import (
	"github.com/gaols/easyssh"
)

// Tar pack the targetPath and put tarball to tgzPath, targetPath and tgzPath should both the absolute path.
func Tar(tgzPath, targetPath string) Command {
	return BuildSimpleCommand(func(deploy Deploy, config *ServerConfig, args ...interface{}) error {
		tgzPath := args[0].(string)
		targetPath := args[1].(string)
		return easyssh.Tar(tgzPath, targetPath)
	}, tgzPath, targetPath)
}

// UnTar unpack the tarball specified by tgzPath and extract it to the path specified by targetPath
func UnTar(tgzPath, targetPath string) Command {
	return BuildSimpleCommand(func(deploy Deploy, config *ServerConfig, args ...interface{}) error {
		tgzPath := args[0].(string)
		targetPath := args[1].(string)
		return easyssh.UnTar(tgzPath, targetPath)
	}, tgzPath, targetPath)
}
