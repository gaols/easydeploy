package easydeploy

import "github.com/gaols/easyssh"

type Command interface {
	Run() error
}

type LocalCommand struct {
	Cmd string
}

type RemoteCommand struct {
	Cmd          string
	ServerConfig *ServerConfig
}

type UploadCommand struct {
	src  string
	dest string
}

func (localCmd *LocalCommand) Run() error {
	if "" == localCmd.Cmd {
		panic("missing local command")
	}

	_, err := easyssh.Local(localCmd.Cmd)
	return err
}

func (remoteCommand *RemoteCommand) Run() error {
	if "" == remoteCommand.Cmd {
		panic("missing remote command")
	}

	ssh := remoteCommand.ServerConfig.MakeSSHConfig()
	_, _, _, err := ssh.Run(remoteCommand.Cmd, -1)
	return err
}
