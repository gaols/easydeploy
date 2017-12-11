package easydeploy

import (
	"github.com/gaols/easyssh"
	"strings"
	"fmt"
)

type Command interface {
	Run(deployCtx Deploy, srvConf *ServerConfig) error
}

type LocalCommand struct {
	CmdStr string
}

type RemoteCommand struct {
	CmdStr string
}

type UploadCommand struct {
	LocalPath  string
	RemotePath string
}

func (localCmd *LocalCommand) Run(deployCtx Deploy, srvConf *ServerConfig) error {
	if isBlank(localCmd.CmdStr) {
		panic("missing local command")
	}

	vOutCommand(srvConf, localCmd.CmdStr, "local")
	err := easyssh.RtLocal(localCmd.CmdStr, func(line string, lineType int8) {
		if deployCtx.isVerbose() {
			vOut(srvConf, line)
		}
	})

	return err
}

func (remoteCmd *RemoteCommand) Run(deployCtx Deploy, srvConf *ServerConfig) error {
	if isBlank(remoteCmd.CmdStr) {
		panic("missing remote command")
	}
	vOutCommand(srvConf, remoteCmd.CmdStr, "remote")

	ssh := srvConf.MakeSSHConfig()
	_, err := ssh.RtRun(remoteCmd.CmdStr, func(line string, lineType int) {
		if deployCtx.isVerbose() {
			vOut(srvConf, line)
		}
	}, -1)
	return err
}

func (uploadCmd *UploadCommand) Run(deployCtx Deploy, srvConf *ServerConfig) error {
	if isBlank(uploadCmd.LocalPath) {
		panic("missing local path for uploading")
	}

	if isBlank(uploadCmd.RemotePath) {
		panic("missing remote path for uploading")
	}
	uploadM := fmt.Sprintf("%s -> %s", uploadCmd.LocalPath, uploadCmd.RemotePath)
	vOutCommand(srvConf, uploadM, "upload")

	ssh := srvConf.MakeSSHConfig()
	err := ssh.SafeScp(uploadCmd.LocalPath, uploadCmd.RemotePath)
	if err == nil && deployCtx.isVerbose() {
		vOut(srvConf, fmt.Sprintf("%s upload ok", uploadM))
	}
	return err
}

func isBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func vOut(srvConf *ServerConfig, output string) {
	fmt.Printf("[%s] %s\n", srvConf.Simple(), output)
}

func vOutCommand(srvConf *ServerConfig, cmd string, cmdType string) {
	vOut(srvConf, fmt.Sprintf("Run %s command: %s", cmdType, cmd))
}

type SimpleCommand struct {
	Handler func(Deploy, *ServerConfig, ...interface{}) error
	Args    []interface{}
}

func (cmd *SimpleCommand) Run(deployCtx Deploy, srvConf *ServerConfig) error {
	return cmd.Handler(deployCtx, srvConf, cmd.Args...)
}

func BuildSimpleCommand(fn func(Deploy, *ServerConfig, ...interface{}) error, args ...interface{}) Command {
	return &SimpleCommand{
		Handler: fn,
		Args:    args,
	}
}
