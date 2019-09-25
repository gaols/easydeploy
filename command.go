package easydeploy

import (
	"github.com/gaols/easyssh"
	"strings"
	"fmt"
)

// Command is the command interface
type Command interface {
	Run(deployCtx Deploy, srvConf *ServerConfig) error
}

// LocalCommand is a local command
type LocalCommand struct {
	CmdStr string
}

// RemoteCommand is a remote command
type RemoteCommand struct {
	CmdStr string
}

// UploadCommand is the upload command
type UploadCommand struct {
	LocalPath  string
	RemotePath string
}

// DownloadCommand is the download command
type DownloadCommand struct {
	LocalPath  string
	RemotePath string
}

// Run the local command
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

// Run the remote command on the remote server 
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

// Run the download command 
func (uploadCmd *DownloadCommand) Run(deployCtx Deploy, srvConf *ServerConfig) error {
	if isBlank(uploadCmd.LocalPath) {
		panic("missing local path for downloading")
	}

	if isBlank(uploadCmd.RemotePath) {
		panic("missing remote path for downloading")
	}
	downloadM := fmt.Sprintf("%s -> %s", uploadCmd.RemotePath, uploadCmd.LocalPath)
	vOutCommand(srvConf, downloadM, "download")

	ssh := srvConf.MakeSSHConfig()
	err := ssh.DownloadF(uploadCmd.RemotePath, uploadCmd.LocalPath)
	if err == nil && deployCtx.isVerbose() {
		vOut(srvConf, fmt.Sprintf("%s download ok", downloadM))
	}
	return err
}

// Run the upload command 
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

// SimpleCommand represents a simple customized command
type SimpleCommand struct {
	Handler func(Deploy, *ServerConfig, ...interface{}) error
	Args    []interface{}
}

// Run this simple customized command 
func (cmd *SimpleCommand) Run(deployCtx Deploy, srvConf *ServerConfig) error {
	return cmd.Handler(deployCtx, srvConf, cmd.Args...)
}

// BuildSimpleCommand build the simple command
func BuildSimpleCommand(fn func(Deploy, *ServerConfig, ...interface{}) error, args ...interface{}) Command {
	return &SimpleCommand{
		Handler: fn,
		Args:    args,
	}
}
