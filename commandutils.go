package easydeploy

import (
	"path/filepath"
	"github.com/gaols/easyssh"
	"errors"
	"fmt"
)

// Tar pack the targetPath and put tarball to tgzPath, targetPath and tgzPath should both the absolute path.
func Tar(tgzPath, targetPath string) Command {
	return BuildSimpleCommand(func(deploy Deploy, config *ServerConfig, args ...interface{}) error {
		tgzPath, ok1 := args[0].(string)
		targetPath, ok2 := args[1].(string)

		if !ok1 || !ok2 {
			return errors.New(fmt.Sprintf("invalid args: %s", args))
		}

		targetPathDir := filepath.Dir(easyssh.RemoveTrailingSlash(targetPath))
		target := filepath.Base(easyssh.RemoveTrailingSlash(targetPath))
		_, err := easyssh.Local("tar czf %s -C %s %s", tgzPath, targetPathDir, target)
		return err
	}, tgzPath, targetPath)
}

// UnTar unpack the tarball specified by tgzPath and extract it to the path specified by targetPath
func UnTar(tgzPath, targetPath string) Command {
	return BuildSimpleCommand(func(deploy Deploy, config *ServerConfig, args ...interface{}) error {
		tgzPath, ok1 := args[0].(string)
		targetPath, ok2 := args[1].(string)

		if !ok1 || !ok2 {
			return errors.New(fmt.Sprintf("invalid args: %s", args))
		}

		_, err := easyssh.Local("tar xf %s -C %s", tgzPath, targetPath)
		return err
	}, tgzPath, targetPath)
}
