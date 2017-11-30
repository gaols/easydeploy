// Package easydeploy makes your deployment easy
package easydeploy

import (
	"github.com/gaols/easyssh"
	"fmt"
	"strings"
)

type Deploy interface {
	RegisterDeployServer(srvConf *ServerConfig)
	OnceBeforeDeploy(fn func() error)
	OnceDoneDeploy(fn func(deployOk bool) error)
	Local(cmd string, args ...interface{})
	Remote(cmd string, args ...interface{})
	Upload(localPath, remotePath string)
	Start() error
	StartAsync() chan error
	Verbose()
}

type Deployer struct {
	SrvConf      []*ServerConfig
	commands     []Command
	onceDoneFn   func(bool) error
	onceBeforeFn func() error
	verbose      bool
}

type ServerConfig struct {
	User         string
	Server       string
	Key          string
	Port         string
	Password     string
}

type Report struct {
}

// Local register a command to be run on localhost
func (sc *Deployer) Local(cmd string, args ...interface{}) {
	sc.commands = append(sc.commands, &LocalCommand{
		CmdStr:  fmt.Sprintf(cmd, args...),
	})
}

// Remote register a command to be run on remote host
func (sc *Deployer) Remote(cmd string, args ...interface{}) {
	sc.commands = append(sc.commands, &RemoteCommand{
		CmdStr:  fmt.Sprintf(cmd, args...),
	})
}

// Remote register a upload command
func (sc *Deployer) Upload(localPath, remotePath string) {
	sc.commands = append(sc.commands, &UploadCommand{
		LocalPath:  localPath,
		RemotePath: remotePath,
	})
}

// Start will start the deploy process
func (sc *Deployer) Start() error {
	err := sc.onceBeforeFn()
	if err != nil {
		return err
	}

	for _, cmd := range sc.commands {
		err := cmd.Run(nil, nil)
		if err != nil {
			return err
		}
	}
	return sc.onceDoneFn(false)
}

// StartAsync will start the deployment in a go routing
func (sc *Deployer) StartAsync() chan error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- sc.Start()
	}()
	return errCh
}

// OnceBeforeDeploy called only once before deployment
func (sc *Deployer) OnceBeforeDeploy(fn func() error) {
	sc.onceBeforeFn = fn
}

// OnceDoneDeploy called only once when deployment is done even if deployment failed
func (sc *Deployer) OnceDoneDeploy(fn func(deployOk bool) error) {
	sc.onceDoneFn = fn
}

// Verbose should be called before deployment start if you want to see the each command output.
func (sc *Deployer) Verbose() {
	sc.verbose = true
}

func (sc *ServerConfig) MakeSSHConfig() *easyssh.SSHConfig {
	return &easyssh.SSHConfig{
		User:     sc.User,
		Server:   sc.Server,
		Key:      sc.Key,
		Port:     sc.Port,
		Password: sc.Password,
	}
}

func (sc *ServerConfig) String() string {
	return fmt.Sprintf("%s@%s:%s/%s", sc.User, sc.Server, sc.Port, sc.Password)
}

// NewSrvConf parse server config from string with format: user@host:port/password
func NewSrvConf(format string) *ServerConfig {
	semIdx := strings.Index(format, ":")
	atIdx := strings.Index(format, "@")
	slashIdx := strings.Index(format, "/")
	if atIdx == -1 || slashIdx == -1 || semIdx == -1 {
		panic("invalid server config format: " + format)
	}

	return &ServerConfig{
		User:     format[0:atIdx],
		Server:   format[atIdx+1:semIdx],
		Port:     format[semIdx+1:slashIdx],
		Password: format[slashIdx+1:],
	}
}
