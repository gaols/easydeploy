// Package easydeploy makes your deployment easy
package easydeploy

import "github.com/gaols/easyssh"

type Deploy interface {
	OnceBefore(fn func() error)
	OnceDone(fn func(deployOk bool) error)
	Local(cmd string, args ...string)
	Remote(cmd string, args ...string)
	Upload(localPath, remotePath string)
	Start() error
}

type ServerConfig struct {
	User     string
	Server   string
	Key      string
	Port     string
	Password string
	Commands []*Command
}

// Local register a command to be run on localhost
func (sc *ServerConfig) Local(cmd string, args ...string) {
}

// Remote register a command to be run on remote host
func (sc *ServerConfig) Remote(cmd string, args ...string) {
}

// Remote register a upload command
func (sc *ServerConfig) Upload(localPath, remotePath string) {
}

// Start will start the deploy process
func (sc *ServerConfig) Start() error {
	return nil
}

func (sc *ServerConfig) OnceBefore(fn func() error) {
}

func (sc *ServerConfig) OnceDone(fn func() error) {
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
