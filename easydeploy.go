// Package easydeploy makes your deployment easy
package easydeploy

import (
	"fmt"
	"github.com/gaols/easyssh"
	"strings"
	"time"
)

// Deploy is the deploy interface
type Deploy interface {
	RegisterDeployServer(srvConf *ServerConfig)
	OnceBeforeDeploy(fn func() error)
	OnceDoneDeploy(fn func(bool) error)
	Local(cmd string, args ...interface{})
	Remote(cmd string, args ...interface{})
	Upload(localPath, remotePath string)
	DownloadF(remotePath, localPath string)
	AddCommand(command Command)
	MaxConcurrency(num int)
	Start() []*DeployReport
	Verbose(verbose bool)
	isVerbose() bool
}

// DeployReport is the report of deploy process
type DeployReport struct {
	SrvConf  *ServerConfig
	error    error
	Start    time.Time
	Consumed int64
	CmdRuns  int
	Cmds     int
}

func (rp *DeployReport) String() string {
	report := "\n[%s] Deploy Report:\nCommands total: %d\nCommands runs : %d\nDeploy result : %s, Consumed: %ds\n"
	deployResult := "OK"
	if rp.error != nil {
		deployResult = "FAIL"
	}
	return fmt.Sprintf(report, rp.SrvConf.Simple(), rp.Cmds, rp.CmdRuns, deployResult, rp.Consumed)
}

// Deployer is the main entry point
type Deployer struct {
	SrvConf       []*ServerConfig
	commands      []Command
	onceDoneFn    func(bool) error
	onceBeforeFn  func() error
	verbose       bool
	verboseCalled bool
	concurrency   int
}

// ServerConfig is the ssh server config used to connect to remote server
type ServerConfig struct {
	User     string
	Server   string
	Key      string
	Port     string
	Password string
}

// Local register a command to be run on localhost
func (sc *Deployer) Local(cmd string, args ...interface{}) {
	sc._local(false, cmd, args...)
}

// SLocal register a command to be run on localhost
func (sc *Deployer) SLocal(cmd string, args ...interface{}) {
	sc._local(true, cmd, args...)
}

func (sc *Deployer) _local(sensitive bool, cmd string, args ...interface{}) {
	localCmd := &LocalCommand{
		CmdStr: fmt.Sprintf(cmd, args...),
	}
	if sensitive {
		localCmd.Sensitive()
	}
	sc.commands = append(sc.commands, localCmd)
}

// RegisterDeployServer register a deploy server
func (sc *Deployer) RegisterDeployServer(srvConf *ServerConfig) {
	sc.SrvConf = append(sc.SrvConf, srvConf)
}

// Remote register a command to be run on remote host
func (sc *Deployer) Remote(cmd string, args ...interface{}) {
	sc._remote(false, cmd, args...)
}

// Remote register a command to be run on remote host
func (sc *Deployer) SRemote(cmd string, args ...interface{}) {
	sc._remote(true, cmd, args...)
}

// Remote register a command to be run on remote host
func (sc *Deployer) _remote(sensitive bool, cmd string, args ...interface{}) {
	remoteCmd := &RemoteCommand{
		CmdStr: fmt.Sprintf(cmd, args...),
	}
	if sensitive {
		remoteCmd.Sensitive()
	}
	sc.commands = append(sc.commands, remoteCmd)
}

// Upload register a upload command
func (sc *Deployer) Upload(localPath, remotePath string) {
	sc.commands = append(sc.commands, &UploadCommand{
		LocalPath:  localPath,
		RemotePath: remotePath,
	})
}

// DownloadF register a upload command
func (sc *Deployer) DownloadF(localPath, remotePath string) {
	sc.commands = append(sc.commands, &DownloadCommand{
		RemotePath: remotePath,
		LocalPath:  localPath,
	})
}

// AddCommand register a custom command
func (sc *Deployer) AddCommand(command Command) {
	sc.commands = append(sc.commands, command)
}

// Start will start the deploy process
func (sc *Deployer) Start() []*DeployReport {
	if !sc.verboseCalled {
		// default verbose
		sc.verbose = true
	}

	if sc.onceBeforeFn != nil {
		err := sc.onceBeforeFn()
		if err != nil {
			_ = sc.onceDoneFn(false)
			return nil
		}
	}

	ret := make([]*DeployReport, 0, len(sc.SrvConf))
	if sc.concurrency == 0 || sc.concurrency > len(sc.SrvConf) {
		sc.concurrency = len(sc.SrvConf)
	}

	reportChan := make(chan *DeployReport, len(sc.SrvConf))
	ctrlChan := make(chan int8, sc.concurrency)
	for _, s := range sc.SrvConf {
		startDeployment(sc, s, sc.commands, reportChan, ctrlChan)
	}
L:
	for {
		select {
		case rp := <-reportChan:
			fmt.Print(rp)
			ret = append(ret, rp)
			if len(ret) == len(sc.SrvConf) {
				break L
			}
		}
	}

	if sc.onceDoneFn != nil {
		if err := sc.onceDoneFn(true); err != nil {
			fmt.Println("once done error: ", err.Error())
		}
	}
	return ret
}

func startDeployment(sc *Deployer, srvConf *ServerConfig, Cmds []Command, reportChan chan *DeployReport, ctrlChan chan int8) {
	rp := &DeployReport{
		SrvConf: srvConf,
		Cmds:    len(Cmds),
	}
	go func() {
		ctrlChan <- 0
		defer func() {
			<-ctrlChan
		}()

		rp.Start = time.Now()
		for i, cmd := range Cmds {
			err := cmd.Run(sc, srvConf)
			rp.CmdRuns = i + 1
			if err != nil {
				rp.error = err
				break
			}
		}
		rp.Consumed = time.Now().Unix() - rp.Start.Unix()
		reportChan <- rp
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
func (sc *Deployer) Verbose(verbose bool) {
	sc.verboseCalled = true
	sc.verbose = verbose
}

// Verbose should be called before deployment start if you want to see the each command output.
func (sc *Deployer) isVerbose() bool {
	return sc.verbose
}

// MaxConcurrency controls how many deployments will run simultaneously if you deploy to multiple servers.
// 0 means all deployments will start asynchronously.
func (sc *Deployer) MaxConcurrency(num int) {
	sc.concurrency = num
}

// MakeSSHConfig converts ServerConfig to easyssh.SSHConfig to ease using easyssh
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

// Simple is too simple too naive to say anything about it
func (sc *ServerConfig) Simple() string {
	return fmt.Sprintf("%s", sc.Server)
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
		Server:   format[atIdx+1 : semIdx],
		Port:     format[semIdx+1 : slashIdx],
		Password: format[slashIdx+1:],
	}
}
