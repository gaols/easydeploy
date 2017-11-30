package easydeploy

import (
	"github.com/gaols/easyssh"
	"strings"
)

type Command interface {
	Run(deployCtx *Deploy, srvConf *ServerConfig) error
}

type LocalCommand struct {
	CmdStr  string
	SrvConf *ServerConfig
}

type RemoteCommand struct {
	CmdStr  string
}

type UploadCommand struct {
	LocalPath  string
	RemotePath string
}

func (localCmd *LocalCommand) Run(deployCtx *Deploy, srvConf *ServerConfig) error {
	if isBlank(localCmd.CmdStr) {
		panic("missing local command")
	}

	_, err := easyssh.Local(localCmd.CmdStr)
	return err
}

func (remoteCommand *RemoteCommand) Run(deployCtx *Deploy, srvConf *ServerConfig) error {
	if isBlank(remoteCommand.CmdStr) {
		panic("missing remote command")
	}

	ssh := srvConf.MakeSSHConfig()
	_, _, _, err := ssh.Run(remoteCommand.CmdStr, -1)
	return err
}

func (uploadCommand *UploadCommand) Run(deployCtx *Deploy, srvConf *ServerConfig) error {
	if isBlank(uploadCommand.LocalPath) {
		panic("missing local path for uploading")
	}

	if isBlank(uploadCommand.RemotePath) {
		panic("missing remote path for uploading")
	}
	ssh := srvConf.MakeSSHConfig()
	return ssh.Scp(uploadCommand.LocalPath, uploadCommand.RemotePath)
}

func isBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
