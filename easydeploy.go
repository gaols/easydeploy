// Package easydeploy makes your deployment easy
package easydeploy

import (
	"github.com/gaols/easyssh"
	"fmt"
	"strings"
	"time"
)

type Deploy interface {
	RegisterDeployServer(srvConf *ServerConfig)
	OnceBeforeDeploy(fn func() error)
	OnceDoneDeploy(fn func(bool) error)
	Local(cmd string, args ...interface{})
	Remote(cmd string, args ...interface{})
	Upload(localPath, remotePath string)
	Start() []*DeployReport
	Verbose()
	isVerbose() bool
}

type DeployReport struct {
	SrvConf  *ServerConfig
	error    error
	Start    time.Time
	Consumed int64
	CmdRuns  int
}

type Deployer struct {
	SrvConf      []*ServerConfig
	commands     []Command
	onceDoneFn   func(bool) error
	onceBeforeFn func() error
	verbose      bool
}

// readonly
type ServerConfig struct {
	User     string
	Server   string
	Key      string
	Port     string
	Password string
}

// Local register a command to be run on localhost
func (sc *Deployer) Local(cmd string, args ...interface{}) {
	sc.commands = append(sc.commands, &LocalCommand{
		CmdStr: fmt.Sprintf(cmd, args...),
	})
}

// Local register a command to be run on localhost
func (sc *Deployer) RegisterDeployServer(srvConf *ServerConfig) {
	sc.SrvConf = append(sc.SrvConf, srvConf)
}

// Remote register a command to be run on remote host
func (sc *Deployer) Remote(cmd string, args ...interface{}) {
	sc.commands = append(sc.commands, &RemoteCommand{
		CmdStr: fmt.Sprintf(cmd, args...),
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
func (sc *Deployer) Start() []*DeployReport {
	reportChan := make(chan *DeployReport, len(sc.SrvConf))
	ret := make([]*DeployReport, 0, len(sc.SrvConf))
	err := sc.onceBeforeFn()
	if err != nil {
		sc.onceDoneFn(false)
		return nil
	}

	for _, s := range sc.SrvConf {
		startDeployment(sc, s, sc.commands, reportChan)
	}
L:
	for {
		select {
		case rp := <-reportChan:
			ret = append(ret, rp)
			if len(ret) == len(sc.SrvConf) {
				break L
			}
		}
	}

	if err := sc.onceDoneFn(true); err != nil {
		fmt.Println("once done error: ", err.Error())
	}
	return ret
}

func startDeployment(sc *Deployer, srvConf *ServerConfig, Commands []Command, reportChan chan *DeployReport) {
	rp := &DeployReport{
		SrvConf: srvConf,
	}
	go func() {
		rp.Start = time.Now()
		for i, cmd := range Commands {
			err := cmd.Run(sc, srvConf)
			rp.Consumed = time.Now().Unix() - rp.Start.Unix()
			if err != nil {
				rp.error = err
				break
			}
			rp.CmdRuns = i + 1
			reportChan <- rp
		}
	}()
}

// OnceBeforeDeploy called only once before deployment
func (sc *Deployer) OnceBeforeDeploy(fn func() error) {
	sc.onceBeforeFn = fn
}

// OnceDoneDeploy called only once when deployment is done even if deployment failed
func (sc *Deployer) OnceDoneDeploy(fn func(bool) error) {
	sc.onceDoneFn = fn
}

// Verbose should be called before deployment start if you want to see the each command output.
func (sc *Deployer) Verbose() {
	sc.verbose = true
}

// Verbose should be called before deployment start if you want to see the each command output.
func (sc *Deployer) isVerbose() bool {
	return sc.verbose
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

func (sc *ServerConfig) Simple() string {
	return fmt.Sprintf("%s@%s", sc.User, sc.Server)
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
